package shell

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"github.com/Yulian302/qugopy/internal/tasks"
	"github.com/Yulian302/qugopy/internal/trie"
	"github.com/Yulian302/qugopy/logging"
	"github.com/Yulian302/qugopy/models"
	"github.com/Yulian302/qugopy/shell/internal"
	"github.com/go-redis/redis"
)

type Shell struct {
	input                  []byte
	wordInput              []byte
	runeSuggestions        []string
	isChangedInput         bool
	currGroup              int
	lastSuggestionsPrinted int
	cursorPos              int

	tokenTrie *trie.TrieToken
	runeTrie  *trie.TrieRune
	history   *internal.RingBuffer
}

func NewShell() *Shell {
	return &Shell{
		input:          make([]byte, 0, 256),
		wordInput:      make([]byte, 0, 64),
		tokenTrie:      trie.NewTokenTrie(),
		runeTrie:       trie.NewRuneTrie(),
		history:        internal.NewRingBuffer(50),
		isChangedInput: true,
		currGroup:      0,
		cursorPos:      0,
	}
}

func (sh *Shell) EraseCharacter(stdout bool) {
	if sh.cursorPos == 0 {
		return
	}
	sh.input = append(sh.input[:sh.cursorPos-1], sh.input[sh.cursorPos:]...)
	sh.cursorPos--

	if stdout {
		sh.redrawInput()
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
		if next, ok := curr.Children[token]; ok {
			curr = next
			continue
		}

		if next, ok := curr.Children["*"]; ok {
			curr = next
			continue
		}

		return nil
	}
	result := make([]string, 0, len(curr.Children))
	for child := range curr.Children {
		result = append(result, child)
	}
	return result
}

func (sh *Shell) getWordUnderCursor() []byte {
	if len(sh.input) == 0 || sh.cursorPos > len(sh.input) {
		return nil
	}

	start := sh.cursorPos
	end := sh.cursorPos

	// Go left to find start of word
	for start > 0 && sh.input[start-1] != ' ' {
		start--
	}

	// Go right to find end of word
	for end < len(sh.input) && sh.input[end] != ' ' {
		end++
	}

	return sh.input[start:end]
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
	sh.wordInput = sh.getWordUnderCursor()
	sh.isChangedInput = true
	sh.eraseSuggestions(sh.lastSuggestionsPrinted)
}

func (sh *Shell) handleEraseWord() {
	for len(sh.input) > 0 && sh.input[len(sh.input)-1] != SPACE {
		sh.EraseCharacter(true)
		sh.wordInput = sh.getWordUnderCursor()
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

func (sh *Shell) redrawInput() {
	// Move cursor to start of line
	fmt.Print("\rqugopy> ")

	// Print full input
	os.Stdout.Write(sh.input)

	// Clear any leftovers if shrinking
	fmt.Print("\033[K")

	// Move cursor back to proper position
	back := len(sh.input) - sh.cursorPos
	if back > 0 {
		fmt.Printf("\033[%dD", back) // move left N times
	}
}
func (sh *Shell) handleAppendChar(b byte) {
	if sh.cursorPos > len(sh.input) {
		sh.cursorPos = len(sh.input)
	}

	sh.input = append(sh.input[:sh.cursorPos], append([]byte{b}, sh.input[sh.cursorPos:]...)...)
	sh.cursorPos++

	sh.wordInput = sh.getWordUnderCursor()
	sh.isChangedInput = true

	sh.eraseSuggestions(sh.lastSuggestionsPrinted)
	sh.redrawInput()
}

func (sh *Shell) handleShowSuggestions() {
	if !sh.isChangedInput {
		return
	}

	tokens := sh.getInputTokens()
	sh.currGroup = len(tokens)

	if sh.cursorPos < len(sh.input) && sh.input[sh.cursorPos] != ' ' {
		// if inside a word, subtract 1 to keep currGroup accurate
		sh.currGroup--
	}

	sh.runeSuggestions = nil
	if len(tokens) == 0 {
		sh.runeSuggestions = sh.runeTrie.GetAllWords(1)
	} else if len(sh.wordInput) > 0 {
		word := string(sh.wordInput)
		if strings.ContainsAny(word, "*?") {
			sh.runeSuggestions = sh.runeTrie.FuzzySearch(word, sh.currGroup)
		} else {
			sh.runeSuggestions = sh.runeTrie.SearchPrefix(word, true, sh.currGroup)
		}
	} else {
		nextTokens := sh.getNextTokensFromTokenTrie(tokens)
		sh.runeSuggestions = append(sh.runeSuggestions, nextTokens...)
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

func splitCommandLine(input string) []string {
	var args []string
	var buf strings.Builder
	var inQuote rune
	escaped := false

	for _, r := range input {
		switch {
		case escaped:
			buf.WriteRune(r)
			escaped = false
		case r == '\\':
			escaped = true
		case inQuote != 0:
			if r == inQuote {
				inQuote = 0
			} else {
				buf.WriteRune(r)
			}
		case r == '"' || r == '\'':
			inQuote = r
		case r == ' ' || r == '\t':
			if buf.Len() > 0 {
				args = append(args, buf.String())
				buf.Reset()
			}
		default:
			buf.WriteRune(r)
		}
	}

	if buf.Len() > 0 {
		args = append(args, buf.String())
	}
	return args
}

func parseArgs(line string) map[string]string {
	args := map[string]string{}
	tokens := splitCommandLine(line)

	for i := 0; i < len(tokens); i++ {
		if strings.HasPrefix(tokens[i], "--") {
			key := strings.TrimPrefix(tokens[i], "--")
			if i+1 < len(tokens) && !strings.HasPrefix(tokens[i+1], "--") {
				args[key] = tokens[i+1]
				i++
			} else {
				args[key] = "" // handle flag without value
			}
		}
	}
	return args
}

func parseTaskFromCmd(line string) (models.Task, error) {
	args := parseArgs(line)

	task := models.Task{}
	taskValue := reflect.ValueOf(&task).Elem()
	taskType := taskValue.Type()

	for i := 0; i < taskValue.NumField(); i++ {
		field := taskType.Field(i)
		jsonTag := field.Tag.Get("json")
		if !field.IsExported() {
			continue
		}

		rawValue, ok := args[jsonTag]
		if !ok {
			continue
		}

		fieldValue := taskValue.Field(i)
		fieldType := field.Type

		switch fieldType.Kind() {
		case reflect.String:
			fieldValue.SetString(rawValue)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if intVal, err := strconv.ParseInt(rawValue, 10, 64); err == nil {
				fieldValue.SetInt(intVal)
			} else {
				return models.Task{}, fmt.Errorf("⚠️ Int parse error for %s: %v\n", jsonTag, err)
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if uintVal, err := strconv.ParseUint(rawValue, 10, 64); err == nil {
				fieldValue.SetUint(uintVal)
			} else {
				return models.Task{}, fmt.Errorf("⚠️ Uint parse error for %s: %v\n", jsonTag, err)
			}

		default:
			if fieldType == reflect.TypeOf(json.RawMessage{}) {
				fieldValue.Set(reflect.ValueOf(json.RawMessage(rawValue)))
			} else {
				return models.Task{}, fmt.Errorf("⚠️ Unknown field type for %s: %v\n", jsonTag, fieldType)
			}
		}
	}

	return task, nil
}

func (sh *Shell) Start(tokenGroups [][]string, rdb *redis.Client) {
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
			case byte(ESC):
				escSeq := make([]byte, 2)
				_, err := os.Stdin.Read(escSeq)
				if err != nil {
					continue
				}
				switch string(append([]byte{0x1b}, escSeq...)) {
				// left
				case "\x1b[D":
					if sh.cursorPos > 0 {
						sh.cursorPos--
						fmt.Print("\033[D")
					}
				// right
				case "\x1b[C":
					if sh.cursorPos < len(sh.input) {
						sh.cursorPos++
						fmt.Print("\033[C")
					}
				// up
				case "\x1b[A":
					if cmd, ok := sh.history.Prev(); ok {
						sh.input = []byte(cmd)
						sh.cursorPos = len(sh.input)
						sh.redrawInput()
					}
				// down
				case "\x1b[B":
					if cmd, ok := sh.history.Next(); ok {
						sh.input = []byte(cmd)
						sh.cursorPos = len(sh.input)
						sh.redrawInput()
					} else {
						sh.input = sh.input[:0]
						sh.cursorPos = 0
						sh.redrawInput()
					}
				// clear screen
				case "\x1bc":
					fmt.Print(CLEAR_SCREEN)
					sh.redrawInput()
					continue
				}

			default:
				sh.handleAppendChar(b)
			}
		}

		line := string(sh.input)
		fmt.Println()
		if strings.TrimSpace(line) == "exit" {
			fmt.Println("Goodbye...")
			break
		}
		sh.history.Add(line)

		task, err := parseTaskFromCmd(line)
		if err != nil {
			fmt.Println("Could not process task!")
			fmt.Printf("Error: %v\n", err)
		}
		if err := tasks.EnqueueTask(task, rdb); err != nil {
			logging.DebugLog(fmt.Sprintf("task could not be added: %v", err))

			if len(sh.input) == 0 {
				fmt.Println("(empty)")
			}
			fmt.Printf("Invalid command: %s\n", line)
			continue
		}
		fmt.Println("Task added successfully!")
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

func StartInteractiveShell(rdb *redis.Client) {
	sh := NewShell()
	sh.Start(tokenGroups, rdb)
	os.Exit(0)
}
