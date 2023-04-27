package expressions

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrorUnterminated    = errors.New("non-terminated statement in expression")
	ErrorEmptyStatement  = errors.New("empty statement in expression")
	ErrorMissingFunction = errors.New("missing function")
)

type DetailedError struct {
	Err     error
	Context string
	Index   int
}

func (s *DetailedError) Error() string {
	return fmt.Sprintf("At `%s` (%d): %v", s.Context, s.Index, s.Err)
}

func (s *DetailedError) Unwrap() error {
	return s.Err
}

type CompilerErrors struct {
	Errors     []*DetailedError
	Expression string
}

func (s *CompilerErrors) Error() string {
	if len(s.Errors) == 1 {
		return s.Errors[0].Error()
	}
	var sb strings.Builder
	sb.WriteString("Compiler Errors in: `")
	sb.WriteString(s.Expression)
	sb.WriteString("`\n")
	for _, e := range s.Errors {
		sb.WriteString("  ")
		sb.WriteString(e.Error())
		sb.WriteString("\n")
	}
	return sb.String()
}

func (s *CompilerErrors) Unwrap() error {
	if len(s.Errors) > 0 {
		return s.Errors[0]
	}
	return nil
}

func (s *CompilerErrors) add(underlying error, context string, offset int) {
	s.Errors = append(s.Errors, &DetailedError{underlying, context, offset})
}

func (s *CompilerErrors) empty() bool {
	return len(s.Errors) == 0
}

// Inherit all errors from another compiler error set, and offset the index eg. if nested compile
func (s *CompilerErrors) inherit(other *CompilerErrors, offset int) {
	for _, oe := range other.Errors {
		s.add(oe.Err, oe.Context, oe.Index+offset)
	}
}
