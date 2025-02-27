//go:build linux && cgo && pcre2

package fastregex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaptureGroupCount(t *testing.T) {
	re, err := Compile("this is (.+) test (.+)")
	assert.NoError(t, err)
	assert.Equal(t, 3, re.CreateInstance().(*pcre2Regexp).GroupCount())
}

// pcre1: 500ns
// pcre1-jit: 400ns
// pcre2: 542ns
// pcre2-jit: 422ns / 310 ns / 237ns
func BenchmarkPCRESubMatch(b *testing.B) {
	re := MustCompile(`t(\w+)`).CreateInstance()
	b := []byte("this is a test")
	for i := 0; i < b.N; i++ {
		re.FindSubmatchIndex(b)
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
