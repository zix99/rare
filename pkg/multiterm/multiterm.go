package multiterm

import (
	"fmt"
	"os"
)

type TermWriter struct {
	cursor       int
	cursorHidden bool
	maxLine      int

	ClearLine  bool
	HideCursor bool
}

var _ MultilineTerm = &TermWriter{}

func New() *TermWriter {
	return &TermWriter{
		cursor:     0,
		maxLine:    0,
		ClearLine:  true,
		HideCursor: true,
	}
}

func (s *TermWriter) WriteForLinef(line int, format string, args ...interface{}) {
	s.WriteForLine(line, fmt.Sprintf(format, args...))
}

func (s *TermWriter) WriteForLine(line int, text string) {
	if s.HideCursor && !s.cursorHidden {
		hideCursor()
		s.cursorHidden = true
	}

	s.goTo(line)
	s.writeAtCursor(text)
}

func (s *TermWriter) Close() {
	s.goTo(s.maxLine)
	fmt.Println() // Put cursor after last line
	if s.cursorHidden {
		showCursor()
	}
}

func (s *TermWriter) goTo(line int) {
	if line > s.maxLine {
		s.maxLine = line
	}
	for i := s.cursor; i < line; i++ {
		fmt.Print("\n")
		s.cursor++
	}
	for i := s.cursor; i > line; i-- {
		moveUp(1)
		s.cursor--
	}

	fmt.Print("\r")
}

func (s *TermWriter) writeAtCursor(text string) {
	WriteLineNoWrap(os.Stdout, text)
	if s.ClearLine {
		eraseRemainingLine()
	}
}
