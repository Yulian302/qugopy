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

type Shell struct {
	input                  []byte
	wordInput              []byte
	runeSuggestions        []string
	isChangedInput         bool
	currGroup              int
	lastSuggestionsPrinted int

	tokenTrie *trie.TrieToken
	runeTrie  *trie.TrieRune
}

func NewShell() *Shell {
	return &Shell{
		input:          make([]byte, 0, 256),
		wordInput:      make([]byte, 0, 64),
		tokenTrie:      trie.NewTokenTrie(),
		runeTrie:       trie.NewRuneTrie(),
		isChangedInput: true,
		currGroup:      0,
	}
}

func (sh *Shell) EraseCharacter(stdout bool) {
	if len(sh.input) == 0 {
		return
	}
	sh.input = sh.input[:len(sh.input)-1]
	if stdout {
		os.Stdout.Write(ERASE_CHAR)
	}
}

func (sh *Shell) printSuggestions(suggestions []string) {
	if len(suggestions) == 0 {
		return
	}
	sh.eraseSuggestions(sh.lastSuggestionsPrinted)

	os.Stdout.Write(SAVE_CURSOR_POS)
	os.Stdout.Write(MOVE_CURSOR_DOWN_LEFT)
	os.Stdout.Write(DIM_TEXT)

	for _, s := range suggestions {
		os.Stdout.Write([]byte(s + "\n"))
	}

	os.Stdout.Write(RESET_ALL_MODES)
	os.Stdout.Write(RESTORE_CURSOR_POS)
	os.Stdout.Sync()
}

func (sh *Shell) eraseSuggestions(n int) {
	if n == 0 {
		return
	}
	os.Stdout.Write(SAVE_CURSOR_POS)
	os.Stdout.Write(MOVE_CURSOR_DOWN_LEFT)
	for i := 0; i < n; i++ {
		os.Stdout.Write(ERASE_ENTIRE_LINE)
		if i != n-1 {
			os.Stdout.Write([]byte(MOVE_CURSOR_DOWN_LEFT))
		}
	}
	os.Stdout.Write([]byte(fmt.Sprintf(MOVE_CURSOR_PREV_N_BEG, n)))
	os.Stdout.Write(RESTORE_CURSOR_POS)
	os.Stdout.Sync()
}

func (sh *Shell) getInputTokens() []string {
	return strings.Fields(string(sh.input))
}

func (sh *Shell) getNextTokensFromTokenTrie(tokens []string) []string {
	curr := sh.tokenTrie.Root
	for _, token := range tokens {
		next, ok := curr.Children[token]
		if !ok {
			return nil
		}
		curr = next
	}
	result := make([]string, 0, len(curr.Children))
	for child := range curr.Children {
		result = append(result, child)
	}
	return result
}

func (sh *Shell) getLastWord() []byte {
	words := strings.Fields(string(sh.input))
	if len(words) == 0 {
		return nil
	}
	return []byte(words[len(words)-1])
}

func (sh *Shell) populateRuneTrie(tokenGroups [][]string) {
	maxCols := 0
	for _, row := range tokenGroups {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}

	columns := make([][]string, maxCols)
	for i := 0; i < maxCols; i++ {
		columns[i] = []string{}
	}

	for _, row := range tokenGroups {
		for colIdx, val := range row {
			columns[colIdx] = append(columns[colIdx], val)
		}
	}

	for idx, column := range columns {
		sh.runeTrie.Populate(column, idx+1)
	}
}

func (sh *Shell) handleBackspace() {
	sh.EraseCharacter(true)
	sh.wordInput = sh.getLastWord()
	sh.isChangedInput = true
	sh.eraseSuggestions(sh.lastSuggestionsPrinted)
}

func (sh *Shell) handleEraseWord() {
	for len(sh.input) > 0 && sh.input[len(sh.input)-1] != SPACE {
		sh.EraseCharacter(true)
		sh.wordInput = sh.getLastWord()
	}
	for len(sh.input) > 0 && sh.input[len(sh.input)-1] == SPACE {
		sh.EraseCharacter(true)
	}

	sh.eraseSuggestions(sh.lastSuggestionsPrinted)
	sh.lastSuggestionsPrinted = 0
	sh.isChangedInput = true
}

func (sh *Shell) handleEraseAll() {
	for len(sh.input) > 0 {
		sh.EraseCharacter(true)
	}
	sh.wordInput = sh.wordInput[:0]
	sh.eraseSuggestions(sh.lastSuggestionsPrinted)
	sh.lastSuggestionsPrinted = 0
	sh.isChangedInput = true
}

func (sh *Shell) handleAppendChar(b byte, buffer []byte) {
	sh.input = append(sh.input, b)
	sh.wordInput = append(sh.wordInput, b)
	if b == SPACE {
		sh.wordInput = sh.wordInput[:0]
	}
	sh.isChangedInput = true
	sh.eraseSuggestions(sh.lastSuggestionsPrinted)
	os.Stdout.Write(buffer)
}

func (sh *Shell) handleShowSuggestions() {
	if sh.isChangedInput {
		tokens := sh.getInputTokens()
		sh.currGroup = len(tokens)

		sh.runeSuggestions = nil
		if len(tokens) == 0 {
			sh.runeSuggestions = sh.runeTrie.GetAllWords(1)
		} else if len(sh.wordInput) > 0 {
			sh.runeSuggestions = sh.runeTrie.SearchPrefix(string(sh.wordInput), true, sh.currGroup)
		} else {
			nextTokens := sh.getNextTokensFromTokenTrie(tokens)
			for _, allowed := range nextTokens {
				sh.runeSuggestions = append(sh.runeSuggestions, allowed)
			}
		}

		suggestionMap := map[string]struct{}{}
		for _, s := range sh.runeSuggestions {
			suggestionMap[s] = struct{}{}
		}
		allSuggestions := make([]string, 0, len(suggestionMap))
		for s := range suggestionMap {
			allSuggestions = append(allSuggestions, s)
		}

		sh.eraseSuggestions(sh.lastSuggestionsPrinted)
		sh.printSuggestions(allSuggestions)
		sh.lastSuggestionsPrinted = len(allSuggestions)

		sh.isChangedInput = false
	}
}

func (sh *Shell) Start(tokenGroups [][]string) {
	fd := int(os.Stdin.Fd())
	if err := enableTermRawMode(fd); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to enter raw mode:", err)
		os.Exit(1)
	}
	defer disableRawMode(fd)

	sh.tokenTrie.Populate(tokenGroups)
	sh.populateRuneTrie(tokenGroups)

	buffer := make([]byte, 1)

	for {
		fmt.Print("qugopy> ")
		sh.input, sh.wordInput = sh.input[:0], sh.wordInput[:0]
		sh.isChangedInput = true
		sh.currGroup = 0

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
				sh.handleBackspace()
			case CTRL_C:
				fmt.Println("Exiting...")
				return
			case OPTION_BACKSPACE:
				sh.handleEraseWord()
			case HORIZONTAL_TAB:
				sh.handleShowSuggestions()
			case CMD_BACKSPACE:
				sh.handleEraseAll()
			default:
				sh.handleAppendChar(b, buffer)
			}
		}

		line := string(sh.input)
		fmt.Println()
		if strings.TrimSpace(line) == "exit" {
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
		return errors.New("Error disabling raw mode")
	}
	return nil
}

func StartInteractiveShell() {
	sh := NewShell()
	sh.Start(tokenGroups)
}
