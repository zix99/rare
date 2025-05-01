package multiterm

import (
	"fmt"
	"io"
)

type VirtualTerm struct {
	lines  []string
	closed bool
}

var _ MultilineTerm = &VirtualTerm{}

func NewVirtualTerm() *VirtualTerm {
	return NewVirtualTermEx(0, 10)
}

func NewVirtualTermEx(size, cap int) *VirtualTerm {
	return &VirtualTerm{
		lines: make([]string, size, cap),
	}
}

func (s *VirtualTerm) WriteForLine(line int, text string) {
	if s.closed {
		panic("virtualterm closed")
	}

	for line >= len(s.lines) {
		s.lines = append(s.lines, "")
	}

	s.lines[line] = text
}

func (s *VirtualTerm) WriteForLinef(line int, format string, args ...interface{}) {
	s.WriteForLine(line, fmt.Sprintf(format, args...))
}

// Close the virtual term. Doesn't ever really need to close, but useful in testing
func (s *VirtualTerm) Close() {
	s.closed = true
}

// IsClosed checks if Close() was called
func (s *VirtualTerm) IsClosed() bool {
	return s.closed
}

func (s *VirtualTerm) Get(line int) string {
	if line >= len(s.lines) || line < 0 {
		return ""
	}
	return s.lines[line]
}

func (s *VirtualTerm) LineCount() int {
	return len(s.lines)
}

// WriteToOutput writes to a terminal, preventing any potential wrapping
func (s *VirtualTerm) WriteToOutput(out io.StringWriter) {
	for _, line := range s.lines {
		WriteLineNoWrap(out, line)
		out.WriteString("\n")
	}
}
