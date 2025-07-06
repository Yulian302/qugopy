package shell


import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"github.com/Yulian302/qugopy/internal/trie"
)

var (
	originalState          syscall.Termios
	tokenTrie              *trie.TrieToken = trie.NewTokenTrie()
	runeTrie               *trie.TrieRune  = trie.NewRuneTrie()
	lastSuggestionsPrinted int
)

var (
	BACKSPACE_1      uint8 = 127
	BACKSPACE_2      uint8 = 8
	ENTER_1          uint8 = '\r'
	ENTER_2          uint8 = '\n'
	CTRL_C           uint8 = 3
	OPTION_BACKSPACE uint8 = 23
	SPACE            uint8 = 32
	CMD_BACKSPACE    uint8 = 21
	HORIZONTAL_TAB   uint8 = 9
)

func EraseCharacter(input *[]byte, stdout bool) {
	if len(*input) == 0 {
		return
	}
	*input = (*input)[:len(*input)-1]
	if stdout {
		os.Stdout.Write([]byte("\b \b"))
	}
}


func printSuggestions(suggestions *[]string) {
	if len(*suggestions) == 0 {
		return
	}

	eraseSuggestions(len(*suggestions))

	// Save cursor
	os.Stdout.Write([]byte("\033[s"))
	// Move cursor down one line to print suggestions
	os.Stdout.Write([]byte("\033[1E"))
	// Dim style
	os.Stdout.Write([]byte("\033[2m"))

	for _, s := range *suggestions {
		os.Stdout.Write([]byte(s + "\n"))
	}

	// Reset style
	os.Stdout.Write([]byte("\033[0m"))
	// Restore cursor
	os.Stdout.Write([]byte("\033[u"))

	os.Stdout.Sync()
}

func eraseSuggestions(n int) {
	// Save cursor pos
	os.Stdout.Write([]byte{0x1B, '[', 's'})

	// Move down one line (where suggestions start)
	os.Stdout.Write([]byte{0x1B, '[', '1', 'E'})

	// For each suggestion line, clear it
	for i := 0; i < n; i++ {
		// Clear entire line
		os.Stdout.Write([]byte{0x1B, '[', '2', 'K'})
		// Move cursor down one line unless last iteration
		if i != n-1 {
			os.Stdout.Write([]byte{0x1B, '[', '1', 'E'})
		}
	}

	// Move cursor back to the line just above suggestions
	os.Stdout.Write([]byte{0x1B, '[', byte(n), 'F'})

	// Restore cursor pos
	os.Stdout.Write([]byte{0x1B, '[', 'u'})
	os.Stdout.Sync()
}


func getInputTokens(input []byte) []string {
	return strings.Fields(string(input))
}

func getNextTokensFromTokenTrie(tokens []string) []string {
	curr := tokenTrie.Root
	for _, token := range tokens {
		if next, ok := curr.Children[token]; ok {
			curr = next
		} else {
			return []string{}
		}
	}

	result := make([]string, 0, len(curr.Children))
	for child := range curr.Children {
		result = append(result, child)
	}
	return result
}

func main() {
	fd := int(os.Stdin.Fd())
	if err := enableTermRawMode(fd); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to enter raw mode:", err)
		os.Exit(1)
	}
	defer disableRawMode(fd)

	tokens := [][]string{
		{"send", "email"},
		{"send", "file"},
		{"set", "name"},
		{"set", "password"},
		{"run", "script"},
		{"run", "executable"},
		{"start", "command", "workers"},
		{"start", "command", "redis"},
	}
	tokenGroups := []struct {
		tokens []string
		group  int
	}{
		{
			tokens: []string{"send", "set", "run", "start"}, group: 1,
		},
		{
			tokens: []string{"email", "file", "name", "password", "script", "executable", "command"}, group: 2,
		},
		{
			tokens: []string{"workers", "redis"}, group: 3,
		},
	}

	tokenTrie.Populate(tokens)
	for _, groupedWords := range tokenGroups {
		runeTrie.Populate(groupedWords.tokens, groupedWords.group)
	}

	buffer := make([]byte, 1)
	wordInput := make([]byte, 0)
	input := make([]byte, 0)

	var cmdSuggestions []string
	var runeSuggestions []string
	var isChangedInput bool
	var currGroup int

	for {
		os.Stdout.Write([]byte("> "))
		input = input[:0]
		wordInput = wordInput[:0]
		isChangedInput = true
		currGroup = 0

	OuterLoop:
		for {
			_, err := os.Stdin.Read(buffer)
			if err != nil {
				fmt.Println("\nRead error:", err)
				return
			}

			b := buffer[0]

			switch b {
			case ENTER_1, ENTER_2:
				break OuterLoop
			case BACKSPACE_1, BACKSPACE_2:
				EraseCharacter(&input, true)
				if b != SPACE {
					EraseCharacter(&wordInput, false)
				}

				isChangedInput = true

			case CTRL_C:
				{
					fmt.Println("Exiting...")
					return
				}
			case OPTION_BACKSPACE:
				{
					if len(input) > 0 {
						for len(input) > 0 && input[len(input)-1] != SPACE {
							EraseCharacter(&input, true)
							EraseCharacter(&wordInput, false)
						}
						// removing trailing space
						for len(input) > 0 && input[len(input)-1] == SPACE {
							EraseCharacter(&input, true)
						}

						eraseSuggestions(lastSuggestionsPrinted)
						lastSuggestionsPrinted = 0
						isChangedInput = true

					}
				}
			case HORIZONTAL_TAB:
				{

					if isChangedInput {
						tokens := getInputTokens(input)
						currGroup = len(tokens)

						cmdSuggestions = nil
						runeSuggestions = nil

						if len(tokens) == 0 {
							// suggest first token
							runeSuggestions = runeTrie.GetAllWords(1)
						} else if len(wordInput) > 0 {
							// mid-token completion (like 'em' → 'email')
							runeSuggestions = runeTrie.SearchPrefix(string(wordInput), true, currGroup)
						} else {
							// Cursor after space → user might want next token OR full command
							cmdSuggestions = tokenTrie.SearchPrefix(tokens, false)

							nextTokens := getNextTokensFromTokenTrie(tokens)
							for _, candidate := range runeTrie.GetAllWords(currGroup + 1) {
								for _, allowed := range nextTokens {
									if candidate == allowed {
										runeSuggestions = append(runeSuggestions, candidate)
										break
									}
								}
							}
						}

						allSuggestions := append(cmdSuggestions, runeSuggestions...)
						eraseSuggestions(lastSuggestionsPrinted)
						printSuggestions(&allSuggestions)
						lastSuggestionsPrinted = len(allSuggestions)
					}
					isChangedInput = false
				}


			case CMD_BACKSPACE:
				{
					for len(input) > 0 {
						EraseCharacter(&input, true)
						if len(wordInput) > 0 {
							EraseCharacter(&wordInput, false)
						}
					}

					eraseSuggestions(lastSuggestionsPrinted)
					lastSuggestionsPrinted = 0
					isChangedInput = true

				}
			default:
				{
					input = append(input, b)
					wordInput = append(wordInput, b)
					if b == SPACE {
						wordInput = wordInput[:0]
					}
					isChangedInput = true
					eraseSuggestions(len(cmdSuggestions) + len(runeSuggestions))

					os.Stdout.Write(buffer)
				}
			}

		}
		line := string(input)
		fmt.Println()

		if line == "exit" {
			fmt.Println("Goodbye...")
			break
		}
		fmt.Println("You typed:", line)
	}
}

func enableTermRawMode(fd int) error {
	var termios syscall.Termios
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TIOCGETA), uintptr(unsafe.Pointer(&termios)))
	if err != 0 {
		return errors.New("Error calling syscall")
	}

	originalState = termios

	termios.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.ISIG
	termios.Iflag &^= syscall.IXON | syscall.ICRNL
	termios.Cc[syscall.VMIN] = 1
	termios.Cc[syscall.VTIME] = 0

	_, _, err = syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TIOCSETA), uintptr(unsafe.Pointer(&termios)))
	if err != 0 {
		return errors.New("Error enabling raw mode")
	}
	return nil
}

func disableRawMode(fd int) error {
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TIOCSETA), uintptr(unsafe.Pointer(&originalState)))
	if err != 0 {
		return errors.New("Erorr disabling raw mode")
	}
	return nil
}

func StartInteractiveShell() {
	main()
}
