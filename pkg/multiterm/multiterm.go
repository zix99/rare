package multiterm

import (
	"fmt"
	"io"
	"os"
)

type TermWriter struct {
	w termCursor

	cursor       int
	cursorHidden bool
	maxLine      int

	ClearLine  bool
	HideCursor bool
}

var _ MultilineTerm = &TermWriter{}

func NewEx(w io.StringWriter) *TermWriter {
	return &TermWriter{
		w:          termCursor{w},
		cursor:     0,
		maxLine:    0,
		ClearLine:  true,
		HideCursor: true,
	}
}

func New() *TermWriter {
	return NewEx(os.Stdout)
}

func (s *TermWriter) WriteForLinef(line int, format string, args ...interface{}) {
	s.WriteForLine(line, fmt.Sprintf(format, args...))
}

func (s *TermWriter) WriteForLine(line int, text string) {
	if s.HideCursor && !s.cursorHidden {
		s.w.hideCursor()
		s.cursorHidden = true
	}

	s.goTo(line)
	s.writeAtCursor(text)
}

func (s *TermWriter) Close() {
	s.goTo(s.maxLine)
	s.w.WriteString("\n") // Put cursor after last line
	if s.cursorHidden {
		s.w.showCursor()
	}
}

func (s *TermWriter) goTo(line int) {
	if line > s.maxLine {
		s.maxLine = line
	}
	for i := s.cursor; i < line; i++ {
		s.w.WriteString("\n")
		s.cursor++
	}
	for i := s.cursor; i > line; i-- {
		s.w.moveUp(1)
		s.cursor--
	}

	s.w.WriteString("\r")
}

func (s *TermWriter) writeAtCursor(text string) {
	WriteLineNoWrap(os.Stdout, text)
	if s.ClearLine {
		s.w.eraseRemainingLine()
	}
}
