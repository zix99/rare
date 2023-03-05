package grok

import (
	"errors"
	"fmt"
	"rare/pkg/grok/stdpat"
	"strings"
)

type PatternLookup interface {
	Lookup(pattern string) (string, bool)
}

type Grok struct {
	patterns PatternLookup
}

func New() *Grok {
	return NewEx(stdpat.Stdlib())
}

func NewEx(patterns PatternLookup) *Grok {
	return &Grok{
		patterns,
	}
}

func (s *Grok) RewriteGrokPattern(str string) (string, error) {
	var sb strings.Builder
	var sub strings.Builder

	var (
		percentActive = false
		inBraces      = false
	)

	runes := []rune(str)
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if inBraces {
			if r == '}' { // end
				expr, err := s.expressionRegex(sub.String())
				if err != nil {
					return "", err
				}
				rExpr, err := s.RewriteGrokPattern(expr) // Recurse
				if err != nil {
					return "", err
				}
				sb.WriteString(rExpr)

				sub.Reset()
				inBraces = false
				percentActive = false
			} else {
				sub.WriteRune(r)
			}
		} else if percentActive {
			if r == '%' { // Escape
				sb.WriteRune(r)
				percentActive = false
			} else if r == '{' { // Start expression
				inBraces = true
				sub.Reset()
			} else { // Not anything special
				sb.WriteRune('%')
				sb.WriteRune(r)
				percentActive = false
			}
		} else if r == '%' {
			percentActive = true
		} else {
			sb.WriteRune(r)
		}
	}

	if inBraces {
		return "", errors.New("unclosed expression")
	}
	if percentActive {
		sb.WriteRune('%')
	}

	return sb.String(), nil
}

func (s *Grok) expressionRegex(grokExpression string) (string, error) {
	parts := strings.Split(grokExpression, ":")

	pattern, ok := s.patterns.Lookup(parts[0])
	if !ok {
		return "<UKN-GROK>", fmt.Errorf("unknown grok expression: %s", grokExpression)
	}

	if len(parts) == 1 {
		return pattern, nil
	}
	return fmt.Sprintf("(?P<%s>%s)", parts[1], pattern), nil
}
