package expressions

// KeyBuilderContextArray is a simple implementation of context with an array of elements
type KeyBuilderContextArray struct {
	Elements []string
}

// GetMatch implements `KeyBuilderContext`
func (s *KeyBuilderContextArray) GetMatch(idx int) string {
	if idx >= 0 && idx < len(s.Elements) {
		return s.Elements[idx]
	}
	return ""
}
