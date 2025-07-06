package main

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/Yulian302/qugopy/internal/trie"
)

var (
	originalState syscall.Termios
	tree          *trie.Trie = trie.NewTrie()
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

func eraseSuggestions() {
	// save cursor
	os.Stdout.Write([]byte{0x1B, '[', 's'})
	// move down 1 line
	os.Stdout.Write([]byte{0x1B, '[', '1', 'E'})
	os.Stdout.Write([]byte{0x1B, '[', '2', 'K'})
	// reset cursor
	os.Stdout.Write([]byte{0x1B, '[', 'u'})
	os.Stdout.Sync()
}

func main() {
	fd := int(os.Stdin.Fd())
	if err := enableTermRawMode(fd); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to enter raw mode:", err)
		os.Exit(1)
	}
	defer disableRawMode(fd)

	keyWords := []string{
		"start",
		"set",
		"send",
		"run",
	}
	tree.Populate(keyWords)

	buffer := make([]byte, 1)
	wordInput := make([]byte, 0)
	input := make([]byte, 0)
	for {
		os.Stdout.Write([]byte("> "))
		input = input[:0]
		wordInput = wordInput[:0]

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
						eraseSuggestions()
					}
				}
			case HORIZONTAL_TAB:
				{
					var suggestions []string
					if len(wordInput) == 0 {
						suggestions = tree.GetAllWords()
					} else {
						suggestions = tree.SearchPrefix(string(wordInput))
					}

					// save cursor
					os.Stdout.Write([]byte{0x1B, '[', 's'})
					// move down 1 line
					os.Stdout.Write([]byte{0x1B, '[', '1', 'E'})
					// dim style
					os.Stdout.Write([]byte{0x1B, '[', '2', 'm'})
					for _, sugg := range suggestions {
						os.Stdout.Write([]byte(sugg))
						os.Stdout.Write([]byte(" "))
					}
					// reset cursor
					os.Stdout.Write([]byte{0x1B, '[', 'u'})
					os.Stdout.Sync()
				}
			case CMD_BACKSPACE:
				{
					for len(input) > 0 {
						EraseCharacter(&input, true)
						if len(wordInput) > 0 {
							EraseCharacter(&wordInput, false)
						}
					}
					eraseSuggestions()
				}
			default:
				{
					input = append(input, b)
					wordInput = append(wordInput, b)
					if b == SPACE {
						wordInput = wordInput[:0]
					}
					eraseSuggestions()
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
