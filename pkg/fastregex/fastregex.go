package fastregex

type Regexp interface {
	Match(b []byte) bool
	MatchString(str string) bool
	FindSubmatchIndex(b []byte) []int
}
