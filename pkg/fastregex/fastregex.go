package fastregex

// An instance of a shareble compiled regex. Assumed to be safe to share across threads/goroutines
type CompiledRegexp interface {
	CreateInstance() Regexp
}

// Regexp servers as an abstraction interface for regex classes
// and shares the same methods as the re2/regexp implementation
// which allows for easy fallback
type Regexp interface {
	Match(b []byte) bool
	MatchString(str string) bool
	FindSubmatchIndex(b []byte) []int
}

// In addition, the following must be provided
var (
	_ string                                                = Version
	_ func(expr string, posix bool) (CompiledRegexp, error) = Compile
	_ func(expr string, posix bool) CompiledRegexp          = MustCompile
)
