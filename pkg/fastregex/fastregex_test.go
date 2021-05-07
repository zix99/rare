package fastregex

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

// These tests should run regardless of the implementation

func TestCompileSuccess(t *testing.T) {
	_, err := Compile("test.*")
	assert.NoError(t, err)
}

func TestCompileError(t *testing.T) {
	_, err := Compile("test(.*")
	assert.Error(t, err)
}

func TestMatch(t *testing.T) {
	assert.True(t, MustCompile("test").CreateInstance().Match([]byte("this is a test")))
	assert.False(t, MustCompile("test").CreateInstance().Match([]byte("this is a tes")))
}

func TestMatchString(t *testing.T) {
	assert.True(t, MustCompile("test").CreateInstance().MatchString("this is a test"))
	assert.False(t, MustCompile("test").CreateInstance().MatchString("this is a tes"))
}

func TestSubMatch(t *testing.T) {
	re := MustCompile("test").CreateInstance()
	assert.Len(t, re.SubexpNameTable(), 0)
	ret := re.FindSubmatchIndex([]byte("this is a test"))
	assert.Len(t, ret, 2)
	assert.Equal(t, []int{10, 14}, ret)
}

func TestSubMatch2(t *testing.T) {
	re := MustCompile("a (\\w+)").CreateInstance()
	ret := re.FindSubmatchIndex([]byte("this is a test"))
	assert.Len(t, ret, 4)
	assert.Equal(t, []int{8, 14, 10, 14}, ret)
}

func TestMatchUnicodeString(t *testing.T) {
	re := MustCompile("test").CreateInstance()
	assert.True(t, re.MatchString("this is ε ζ η a test ✻"))
	assert.Equal(t, []int{19, 23}, re.FindSubmatchIndex([]byte("this is ε ζ η a test ✻")))
}

func TestMatchUnicodeCharacter(t *testing.T) {
	re := MustCompile("ζ").CreateInstance()
	assert.True(t, re.MatchString("this is ε ζ η a test ✻"))
	assert.Equal(t, []int{11, 13}, re.FindSubmatchIndex([]byte("this is ε ζ η a test ✻")))
}

func TestMatchEmptyArray(t *testing.T) {
	re := MustCompile("test").CreateInstance()
	assert.Nil(t, re.FindSubmatchIndex([]byte{}))
	assert.Len(t, re.SubexpNameTable(), 0)
}

func TestCaptureGroupNames(t *testing.T) {
	re := MustCompile(`(?P<num>\d+) (?P<thing>.+) (\w+)`).CreateInstance()
	table := re.SubexpNameTable()
	assert.Len(t, table, 2)
	assert.Equal(t, 1, table["num"])
	assert.Equal(t, 2, table["thing"])
}

// Benchmarks

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

// 305ns
func BenchmarkRE2Match(b *testing.B) {
	re := regexp.MustCompile(`t(\w+)`)
	for i := 0; i < b.N; i++ {
		re.MatchString("this is a test")
	}
}

// pcre1: 500ns
// pcre1-jit: 400ns
// pcre2: 542ns
// pcre2-jit: 422ns / 310 ns / 237ns
func BenchmarkPCRESubMatch(b *testing.B) {
	re := MustCompile(`t(\w+)`).CreateInstance()
	for i := 0; i < b.N; i++ {
		re.FindSubmatchIndex([]byte("this is a test"))
	}
}

// 520ns
func BenchmarkRE2SubMatch(b *testing.B) {
	re := regexp.MustCompile(`t(\w+)`)
	for i := 0; i < b.N; i++ {
		re.FindSubmatchIndex([]byte("this is a test"))
	}
}
