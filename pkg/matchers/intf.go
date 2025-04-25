package matchers

// A thread-safe compiled matcher that can create instances
type Factory interface {
	CreateInstance() Matcher
}

// A non-thread-safe matcher that can be used to find matches
type Matcher interface {
	FindSubmatchIndex(b []byte) []int
	SubexpNameTable() map[string]int

	// Insert FindSubmatchIndex into a preexisting buffer of size MatchBufSize() (no memory alloc)
	// If dst is nil, or undersized, it will be resized/alloc'd
	FindSubmatchIndexDst(b []byte, dst []int) []int
	MatchBufSize() int
}
