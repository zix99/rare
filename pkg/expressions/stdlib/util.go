package stdlib

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
