package stringSplitter

import (
	"rare/pkg/testutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringSplitter(t *testing.T) {
	s := Splitter{
		S:     "abc\x00efg\x00123\x00",
		Delim: "\x00",
	}
	assert.Equal(t, "abc", s.Next())
	assert.Equal(t, "efg", s.Next())
	assert.Equal(t, "123", s.Next())
	assert.False(t, s.Done())
	assert.Equal(t, "", s.Next())
	assert.True(t, s.Done())
}

func TestSplitterNextOk(t *testing.T) {
	s := Splitter{
		S:     "abc\x00efg",
		Delim: "\x00",
	}
	part0, ok0 := s.NextOk()
	assert.Equal(t, "abc", part0)
	assert.True(t, ok0)

	part1, ok1 := s.NextOk()
	assert.Equal(t, "efg", part1)
	assert.True(t, ok1)

	part2, ok2 := s.NextOk()
	assert.Equal(t, "", part2)
	assert.False(t, ok2)
}

// BenchmarkStringSplit-4   	 4282983	       281.6 ns/op	      64 B/op	       1 allocs/op
func BenchmarkStringSplit(b *testing.B) {
	total := 0
	for n := 0; n < b.N; n++ {
		ele := strings.Split("abc\x00efg\x00123\x00", "\x00")
		total += len(ele)
	}
}

// BenchmarkSplitter-4   	15479449	        81.83 ns/op	       0 B/op	       0 allocs/op
func BenchmarkSplitter(b *testing.B) {
	total := 0
	for n := 0; n < b.N; n++ {
		splitter := Splitter{S: "abc\x00efg\x00123\x00", Delim: "\x00"}
		for !splitter.Done() {
			splitter.Next()
			total++
		}
	}
}

func TestZeroAllocs(t *testing.T) {
	testutil.AssertZeroAlloc(t, BenchmarkSplitter)
}

func BenchmarkSplitterNextOk(b *testing.B) {
	total := 0
	for n := 0; n < b.N; n++ {
		splitter := Splitter{S: "abc\x00efg\x00123\x00", Delim: "\x00"}
		for {
			_, ok := splitter.NextOk()
			if !ok {
				break
			}
			total++
		}
	}
}
