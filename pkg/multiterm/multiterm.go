package multiterm

import (
	"fmt"
	"os"
)

type TermWriter struct {
	cursor       int
	cursorHidden bool
	ClearLine    bool
	HideCursor   bool
}

func New() *TermWriter {
	return &TermWriter{
		cursor:     0,
		ClearLine:  true,
		HideCursor: true,
	}
}

func (s *TermWriter) WriteForLine(line int, format string, args ...interface{}) {
	if s.HideCursor && !s.cursorHidden {
		hideCursor()
		s.cursorHidden = true
	}

	s.goTo(line)
	s.writeAtCursor(format, args...)
}

func (s *TermWriter) goTo(line int) {
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

func (s *TermWriter) writeAtCursor(format string, args ...interface{}) {
	WriteLineNoWrap(os.Stdout, fmt.Sprintf(format, args...))
	if s.ClearLine {
		eraseRemainingLine()
	}
}
