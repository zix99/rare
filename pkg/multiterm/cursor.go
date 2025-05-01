package multiterm

import (
	"fmt"
	"io"
)

//lint:file-ignore U1000 Retain useful functions

type termCursor struct {
	io.StringWriter
}

func (s termCursor) moveCursor(line, col int) {
	s.WriteString(termEscape("[%d;%dH", line, col))
}

func (s termCursor) moveUp(n int) {
	s.WriteString(termEscape("[%dA", n))
}

func (s termCursor) hideCursor() {
	s.WriteString(termEscape("[?25l"))
}

func (s termCursor) showCursor() {
	s.WriteString(termEscape("[?25h"))
}

func (s termCursor) eraseRemainingLine() {
	s.WriteString(termEscape("[0K"))
}

func termEscape(format string, args ...interface{}) string {
	const ESCAPE = "\x1b"
	if len(args) == 0 {
		return ESCAPE + format
	}
	return ESCAPE + fmt.Sprintf(format, args...)
}
