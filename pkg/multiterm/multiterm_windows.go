package multiterm

import (
	"os"
	"golang.org/x/sys/windows"
)

func addFlagToWindowsTerm(file *os.File) error {
	term := windows.Handle(file.Fd())
	var outMode uint32
	if err := windows.GetConsoleMode(term, &outMode); err == nil {
		if err := windows.SetConsoleMode(term, outMode | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func init() {
	addFlagToWindowsTerm(os.Stdout)
	addFlagToWindowsTerm(os.Stderr)
}
