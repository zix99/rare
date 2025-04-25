package matchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleMatcherAndFactory(t *testing.T) {
	matcher := ToFactory(&AlwaysMatch{}) // ToFactory isn't necessary, but will exercise the path
	inst := matcher.CreateInstance()

	assert.Empty(t, inst.SubexpNameTable())

	assert.Equal(t, []int{0, 0}, inst.FindSubmatchIndex([]byte{}))
	assert.Equal(t, []int{0, 2}, inst.FindSubmatchIndex([]byte("hi")))
	assert.Equal(t, []int{0, 2}, inst.FindSubmatchIndexDst([]byte("hi"), nil))
	assert.Equal(t, []int{0, 2}, inst.FindSubmatchIndexDst([]byte("hi"), make([]int, 0, inst.MatchBufSize())))
}
