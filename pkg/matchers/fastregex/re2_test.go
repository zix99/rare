package fastregex

import (
	"regexp"
	"testing"
)

// 305ns
func BenchmarkRE2Match(b *testing.B) {
	re := regexp.MustCompile(`t(\w+)`)
	for i := 0; i < b.N; i++ {
		re.MatchString("this is a test")
	}
}

// 520ns
func BenchmarkRE2SubMatch(b *testing.B) {
	re := regexp.MustCompile(`t(\w+)`)
	for i := 0; i < b.N; i++ {
		re.FindSubmatchIndex([]byte("this is a test"))
	}
}
