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
