package testutil

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

// Simplified regex-check that only supports '*' (multi-char) and '?' (single-char) wildcards in the pattern
func AssertPattern(t *testing.T, str, pattern string) {
	t.Helper()
	if err := matchesPattern(pattern, str); err != nil {
		t.Error(err)
	}
}

func matchesPattern(pattern, str string) error {
	p, err := rewritePatternToRegex(pattern)
	if err != nil {
		return fmt.Errorf("invalid pattern: %w", err)
	}

	if !p.MatchString(str) {
		return fmt.Errorf("'%s' does not match pattern '%s'", str, pattern)
	}

	return nil
}

func rewritePatternToRegex(pattern string) (*regexp.Regexp, error) {
	var sb strings.Builder
	sb.Grow(len(pattern))

	sb.WriteRune('^')
	for _, r := range pattern {
		switch r {
		case '*':
			sb.WriteString("(.*)")
		case '?':
			sb.WriteString(".")
		default:
			sb.WriteString(regexp.QuoteMeta(string(r)))
		}
	}
	sb.WriteRune('$')

	return regexp.Compile(sb.String())
}
