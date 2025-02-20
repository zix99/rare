package fastregex

import (
	"rare/pkg/testutil"
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
	ret := re.FindSubmatchIndexDst([]byte("this is a test"), nil)
	assert.Len(t, ret, 2)
	assert.Equal(t, []int{10, 14}, ret)
}

func TestSubMatch2(t *testing.T) {
	re := MustCompile("a (\\w+)").CreateInstance()
	ret := re.FindSubmatchIndexDst([]byte("this is a test"), nil)
	assert.Len(t, ret, 4)
	assert.Equal(t, []int{8, 14, 10, 14}, ret)
}

func TestMatchUnicodeString(t *testing.T) {
	re := MustCompile("test").CreateInstance()
	assert.True(t, re.MatchString("this is ε ζ η a test ✻"))
	assert.Equal(t, []int{19, 23}, re.FindSubmatchIndexDst([]byte("this is ε ζ η a test ✻"), nil))
}

func TestMatchUnicodeCharacter(t *testing.T) {
	re := MustCompile("ζ").CreateInstance()
	assert.True(t, re.MatchString("this is ε ζ η a test ✻"))
	assert.Equal(t, []int{11, 13}, re.FindSubmatchIndexDst([]byte("this is ε ζ η a test ✻"), nil))
}

func TestMatchEmptyArray(t *testing.T) {
	re := MustCompile("test").CreateInstance()
	assert.Nil(t, re.FindSubmatchIndexDst([]byte{}, nil))
	assert.Len(t, re.SubexpNameTable(), 0)
}

func TestCaptureGroupNames(t *testing.T) {
	re := MustCompile(`(?P<num>\d+) (?P<thing>.+) (\w+)`).CreateInstance()
	table := re.SubexpNameTable()
	assert.Len(t, table, 2)
	assert.Equal(t, 1, table["num"])
	assert.Equal(t, 2, table["thing"])
}

func TestMemoryZeroAllocs(t *testing.T) {
	testutil.AssertZeroAlloc(t, BenchmarkFastRegex)
}

func TestMemoryExpectations(t *testing.T) {
	re := MustCompile(`t(\w+)`).CreateInstance()
	d := []byte("hello there bob")

	t.Run("nil alloc", func(t *testing.T) {
		m := re.FindSubmatchIndexDst(d, nil)
		assert.Equal(t, []int{6, 11, 7, 11}, m)
	})

	t.Run("undersized buf alloc", func(t *testing.T) {
		buf := make([]int, 0, 1)
		m := re.FindSubmatchIndexDst(d, buf)
		assert.Equal(t, []int{6, 11, 7, 11}, m)
		testutil.AssertNotSameMemory(t, m, buf)
	})

	t.Run("sized buf alloc", func(t *testing.T) {
		buf := make([]int, 0, re.MatchBufSize())
		m := re.FindSubmatchIndexDst(d, buf)
		assert.Equal(t, []int{6, 11, 7, 11}, m)
		assert.Equal(t, m, buf[:len(m)])
		testutil.AssertSameMemory(t, m, buf)
	})

	t.Run("pre-allocd", func(t *testing.T) {
		buf := make([]int, 2)
		m := re.FindSubmatchIndexDst(d, buf)
		assert.Equal(t, []int{0, 0, 6, 11, 7, 11}, m)
		testutil.AssertNotSameMemory(t, m, buf)
	})
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

// BenchmarkRE2Match-4   	 4593924	       255.4 ns/op	       0 B/op	       0 allocs/op
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
		re.FindSubmatchIndexDst([]byte("this is a test"), nil)
	}
}

// 520ns
func BenchmarkRE2SubMatch(b *testing.B) {
	re := regexp.MustCompile(`t(\w+)`)
	for i := 0; i < b.N; i++ {
		re.FindSubmatchIndex([]byte("this is a test"))
	}
}

// re2-findsubmatch:  BenchmarkFastRegex-4   	 2869810	       401.4 ns/op	      32 B/op	       1 allocs/op
// re2-unsafeExecute: BenchmarkFastRegex-4   	 4154098	       292.3 ns/op	       0 B/op	       0 allocs/op
// pcre2-baseline:    BenchmarkPCRESubMatch-4   	 4560806	       271.8 ns/op	      32 B/op	       0 allocs/op
// pcre2-nopool:      BenchmarkFastRegex-4   	 5465024	       206.9 ns/op	       0 B/op	       0 allocs/op
func BenchmarkFastRegex(b *testing.B) {
	re := MustCompile(`t(\w+)`).CreateInstance()
	d := []byte("this is a test")
	buf := make([]int, 0, re.MatchBufSize())
	for i := 0; i < b.N; i++ {
		re.FindSubmatchIndexDst(d, buf)
	}
}
