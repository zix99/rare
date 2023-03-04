package termstate

import (
	"os"

	"golang.org/x/term"
)

// Returns 'true' if output is being piped (Not char device)
func IsPipedOutput() bool {
	if fi, err := os.Stdout.Stat(); err == nil {
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			return true
		}
	}
	return false
}

// Gets size of the terminal
func GetTermRowsCols() (rows, cols int, ok bool) {
	fd := int(os.Stdout.Fd())
	if !term.IsTerminal(fd) {
		return 0, 0, false
	}

	cols, rows, err := term.GetSize(fd)
	if err != nil {
		return 0, 0, false
	}

	return rows, cols, true
}
