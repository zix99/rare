package matchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleMatcherAndFactory(t *testing.T) {
	matcher := ToFactory(&AlwaysMatch{}) // ToFactory isn't necessary, but will exercise the path
	inst := matcher.CreateInstance()

	assert.Empty(t, inst.SubexpNameTable())

	assert.Equal(t, []int{0, 0}, inst.FindSubmatchIndexDst([]byte{}, nil))
	assert.Equal(t, []int{0, 2}, inst.FindSubmatchIndexDst([]byte("hi"), nil))

	buf := make([]int, 0, inst.MatchBufSize())
	assert.Equal(t, []int{0, 2}, inst.FindSubmatchIndexDst([]byte("hi"), buf))
}

func TestNoAlloc(t *testing.T) {
	b := testing.Benchmark(BenchmarkSimpleMatcher)
	assert.Zero(t, b.AllocedBytesPerOp())
	assert.Zero(t, b.AllocsPerOp())
}

// BenchmarkSimpleMatcher-4   	251675946	         4.761 ns/op	       0 B/op	       0 allocs/op
func BenchmarkSimpleMatcher(b *testing.B) {
	m := ToFactory(&AlwaysMatch{}).CreateInstance()
	d := []byte("hi")
	buf := make([]int, 0, m.MatchBufSize())

	for i := 0; i < b.N; i++ {
		m.FindSubmatchIndexDst(d, buf)
	}
}
