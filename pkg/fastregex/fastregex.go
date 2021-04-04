package fastregex

// An instance of a shareble compiled regex. Assumed to be safe to share across threads/goroutines
type CompiledRegexp interface {
	CreateInstance() Regexp
}

// Regexp serves as an abstraction interface for regex classes
// and shares the same methods as the re2/regexp implementation
// which allows for easy fallback. This interface is expeted
// to only be used by a single thread/goroutine
type Regexp interface {
	Match(b []byte) bool
	MatchString(str string) bool
	FindSubmatchIndex(b []byte) []int
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
