package shell

import "syscall"

var (
	originalState syscall.Termios

	BACKSPACE_1      = uint8(127)
	BACKSPACE_2      = uint8(8)
	ENTER_1          = uint8('\r')
	ENTER_2          = uint8('\n')
	CTRL_C           = uint8(3)
	OPTION_BACKSPACE = uint8(23)
	SPACE            = uint8(32)
	ESC              = 0x1b
	CMD_BACKSPACE    = uint8(21)
	HORIZONTAL_TAB   = uint8(9)

	SAVE_CURSOR_POS       = []byte("\033[s")
	RESTORE_CURSOR_POS    = []byte("\033[u")
	DIM_TEXT              = []byte("\033[2m")
	RESET_ALL_MODES       = []byte("\033[0m")
	MOVE_CURSOR_DOWN_LEFT = []byte("\033[1E")
	ERASE_ENTIRE_LINE     = []byte("\033[2K")
	ERASE_CHAR            = []byte("\b \b")
	CLEAR_SCREEN          = []byte("\033[H\033[2J")
	//moves cursor to beginning of previous line, # lines up
	MOVE_CURSOR_PREV_N_BEG = "\033[%dF"
)
