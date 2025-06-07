package matchers

import (
	"testing"

	"github.com/zix99/rare/pkg/testutil"

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

func TestNoFactory(t *testing.T) {
	matcher := NoFactory(&AlwaysMatch{}).CreateInstance()
	assert.Equal(t, 2, matcher.MatchBufSize())
}

func TestSimpleMemory(t *testing.T) {
	matcher := &AlwaysMatch{}
	buf := make([]int, 0, matcher.MatchBufSize())
	ret := matcher.FindSubmatchIndexDst([]byte("abc"), buf)
	assert.Equal(t, []int{0, 3}, ret)
	testutil.AssertSameMemory(t, buf, ret)
}
