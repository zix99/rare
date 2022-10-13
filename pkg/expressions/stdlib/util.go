package stdlib

import (
	"rare/pkg/expressions"
	"strings"
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

// make a delim-separated array
func makeArray(args ...string) string {
	var sb strings.Builder
	for i := 0; i < len(args); i++ {
		if i > 0 {
			sb.WriteRune(expressions.ArraySeparator)
		}
		sb.WriteString(args[i])
	}
	return sb.String()
}
