package dissect

import "unicode"

// Finds case-insensitive index of second string
// ASSUMES second string is already lowered (optimization)
func indexIgnoreCase(s, loweredSubstr string) int {
	n := len(loweredSubstr)
	switch {
	case n == 0:
		return 0
	case len(s) < n:
		return -1
	case len(s) == n:
		for i := 0; i < n; i++ {
			if unicode.ToLower(rune(s[i])) != rune(loweredSubstr[i]) {
				return -1
			}
		}
		return 0
	default:
		for i := 0; i <= len(s)-n; i++ {
			match := true
			for j := 0; j < n; j++ {
				if unicode.ToLower(rune(s[i+j])) != rune(loweredSubstr[j]) {
					match = false
					break
				}
			}
			if match {
				return i
			}
		}
		return -1
	}

}
