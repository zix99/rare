package stdlib

import (
	"github.com/zix99/rare/pkg/expressions"
)

// Checks if word starts with s
func isPartialString(s, word string) bool {
	if len(s) > len(word) {
		return false
	}

	for i := 0; i < len(s); i++ {
		if s[i] != word[i] {
			return false
		}
	}

	return true
}

// Helper to check if number of arguments is between a min and max
func isArgCountBetween(args []expressions.KeyBuilderStage, min, max int) bool {
	return len(args) >= min && len(args) <= max
}

// Check if string is positive numeric quickly for unix time
// this only works for positive, decimal, non-fractional numbers (eg. just 0-9)
// strconv.ParseInt makes 3 allocs and is significantly slower in the non-numeric case
func simpleParseNumeric(s string) (int64, bool) {
	if len(s) == 0 {
		return 0, false
	}

	var n int64
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0, false
		}
		n = n*10 + int64(c-'0')
	}
	return n, true
}
