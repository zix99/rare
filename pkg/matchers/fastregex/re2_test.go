//go:build !pcre2

package fastregex

import (
	"regexp"
	"testing"

	"github.com/zix99/rare/pkg/testutil"

	"github.com/stretchr/testify/assert"
)

// 305ns
func BenchmarkRE2Match(b *testing.B) {
	re := regexp.MustCompile(`t(\w+)`)
	for i := 0; i < b.N; i++ {
		re.MatchString("this is a test")
	}
}

// BenchmarkRE2SubMatch-4   	 2846446	       431.0 ns/op	      32 B/op	       1 allocs/op
func BenchmarkRE2SubMatch(b *testing.B) {
	re := regexp.MustCompile(`t(\w+)`)
	str := []byte("this is a test")
	for i := 0; i < b.N; i++ {
		re.FindSubmatchIndex(str)
	}
}

func TestMemoryAssumptions(t *testing.T) {
	r := MustCompile(`t(\w+)`)
	str := []byte("this is a test")
	ri := r.CreateInstance()

	// Same memory
	t.Run("same memory", func(t *testing.T) {
		buf := make([]int, 0, ri.MatchBufSize())
		ret := ri.FindSubmatchIndexDst(str, buf)
		assert.Equal(t, []int{0, 4, 1, 4}, ret)
		testutil.AssertSameMemory(t, ret, buf)
	})

	// undersized alloc
	t.Run("undersized", func(t *testing.T) {
		buf := make([]int, 0, 2)
		ret := ri.FindSubmatchIndexDst(str, buf)
		assert.Equal(t, []int{0, 4, 1, 4}, ret)
		testutil.AssertNotSameMemory(t, ret, buf)
	})

	// Nil buf alloc
	t.Run("nilbuf", func(t *testing.T) {
		ret := ri.FindSubmatchIndexDst(str, nil)
		assert.Equal(t, []int{0, 4, 1, 4}, ret)
	})

	// No match
	t.Run("no match", func(t *testing.T) {
		ret := ri.FindSubmatchIndexDst([]byte("none"), nil)
		assert.Nil(t, ret)
	})

	// pad match
	t.Run("pad match", func(t *testing.T) {
		r := MustCompile(`(\d+)(?: (\d+))?`).CreateInstance()
		buf := make([]int, 0, r.MatchBufSize())
		ret := r.FindSubmatchIndexDst([]byte("123"), buf)
		assert.Equal(t, []int{0, 3, 0, 3, -1, -1}, ret)
		testutil.AssertSameMemory(t, buf, ret)
	})
}

func TestDstZeroAlloc(t *testing.T) {
	testutil.AssertZeroAlloc(t, BenchmarkRE2WithDst)
}

// BenchmarkRE2WithDst-4   	 3982665	       295.1 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRE2WithDst(b *testing.B) {
	r := MustCompile(`t(\w+)`)
	str := []byte("this is a test")
	ri := r.CreateInstance()
	buf := make([]int, 0, ri.MatchBufSize())

	for range b.N {
		ri.FindSubmatchIndexDst(str, buf)
	}
}
