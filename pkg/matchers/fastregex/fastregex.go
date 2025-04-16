package fastregex

// An instance of a shareble compiled regex. Assumed to be safe to share across threads/goroutines
type CompiledRegexp interface {
	CreateInstance() Regexp
}

// Regexp serves as an abstraction interface for regex classes
// and shares the same methods as the re2/regexp implementation
// which allows for easy fallback. This interface is expected
// to only be used by a single thread/goroutine
type Regexp interface {
	Match(b []byte) bool
	MatchString(str string) bool

	// Uses buffer specified by dst to fill indexes on match with b
	// Returns nil on no match
	// Dst can be nil, but an alloc will take place. Expects an array
	// with capacity of `MatchBufSize()`, and len of 0
	FindSubmatchIndexDst(b []byte, dst []int) []int

	// Buf size needed to fulfill FindSubmatchIndexDst
	MatchBufSize() int

	// Returns the table of key->match index for FindSubmatch...
	SubexpNameTable() map[string]int
}

// In addition, the following must be provided
var (
	_ string                                                = Version
	_ func(expr string, posix bool) (CompiledRegexp, error) = CompileEx
)

// And lastly, some helpers to bring us closer to regexp module

func Compile(expr string) (CompiledRegexp, error) {
	return CompileEx(expr, false)
}

func MustCompile(expr string) CompiledRegexp {
	re, err := CompileEx(expr, false)
	if err != nil {
		panic(err)
	}
	return re
}
