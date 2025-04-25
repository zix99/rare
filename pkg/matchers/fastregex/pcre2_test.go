//go:build pcre2

package fastregex

import (
	"rare/pkg/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaptureGroupCount(t *testing.T) {
	re, err := Compile("this is (.+) test (.+)")
	assert.NoError(t, err)
	assert.Equal(t, 3, re.CreateInstance().(*pcre2Regexp).GroupCount())
}

func TestPCRESameMemory(t *testing.T) {
	re := MustCompile(`t(\w+)`).CreateInstance()
	sb := []byte("this is a test")
	buf := make([]int, 0, re.MatchBufSize())

	ret := re.FindSubmatchIndexDst(sb, buf)
	testutil.AssertSameMemory(t, buf, ret)
}

func TestPCREZeroAlloc(t *testing.T) {
	testutil.AssertZeroAlloc(t, BenchmarkPCREDst)
}

// pcre1: 500ns
// pcre1-jit: 400ns
// pcre2: 542ns
// pcre2-jit: 422ns / 310 ns / 237ns
func BenchmarkPCRESubMatch(b *testing.B) {
	re := MustCompile(`t(\w+)`).CreateInstance()
	sb := []byte("this is a test")
	for i := 0; i < b.N; i++ {
		re.FindSubmatchIndex(sb)
	}
}

// pcre1: 273ns
// pcre1-jit: 197ns
// pcre2: 308ns
// pcre2-jit: 180ns
func BenchmarkPCREMatch(b *testing.B) {
	re := MustCompile(`t(\w+)`).CreateInstance()
	for i := 0; i < b.N; i++ {
		re.MatchString("this is a test")
	}
}

// BenchmarkPCREDst-4   	 5823752	       206.7 ns/op	       0 B/op	       0 allocs/op
func BenchmarkPCREDst(b *testing.B) {
	re := MustCompile(`t(\w+)`).CreateInstance()
	sb := []byte("this is a test")
	buf := make([]int, 0, re.MatchBufSize())
	for i := 0; i < b.N; i++ {
		re.FindSubmatchIndexDst(sb, buf)
	}
}
