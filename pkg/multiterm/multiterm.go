package multiterm

import (
	"fmt"
)

type TermWriter struct {
	maxLines int
	cursor   int
}

func New(maxLines int) *TermWriter {
	return &TermWriter{
		maxLines: maxLines,
		cursor:   0,
	}
}

func (s *TermWriter) WriteForLine(line int, format string, args ...interface{}) {
	if line >= s.maxLines {
		return
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
	fmt.Printf(format, args...)
}
