package matchers

type AlwaysMatch struct{}

func (s *AlwaysMatch) CreateInstance() Matcher {
	return s
}

func (s *AlwaysMatch) FindSubmatchIndex(b []byte) []int {
	return []int{0, len(b)}
}

func (s *AlwaysMatch) SubexpNameTable() map[string]int {
	return make(map[string]int)
}
