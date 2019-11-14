package color

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	// Force enable, otherwise its init might not be
	Enabled = true
}

func BenchmarkColorReplacer(b *testing.B) {
	s := "This is a test"
	groups := []int{5, 7, 8, 9}

	var out string
	for n := 0; n < b.N; n++ {
		out = WrapIndices(s, groups)
	}

	fmt.Println(out)
}

func TestWrap(t *testing.T) {
	s := Wrap(Red, "hello")
	assert.Contains(t, s, Red)
	assert.Contains(t, s, "hello")
}

func TestWrapDisabled(t *testing.T) {
	Enabled = false
	s := Wrap(Red, "hello")
	assert.Equal(t, "hello", s)
	assert.NotContains(t, s, Red)
	Enabled = true
}

func TestWrapf(t *testing.T) {
	s := Wrapf(Green, "This is %d", 123)
	assert.Contains(t, s, "This is 123")
	assert.Contains(t, s, Green)
	assert.Contains(t, s, Reset)
}

func TestWrapi(t *testing.T) {
	s := Wrapi(Blue, 123)
	assert.Contains(t, s, Blue)
	assert.Contains(t, s, Reset)
	assert.Contains(t, s, "123")
}

func TestWrapIndicesNoGroups(t *testing.T) {
	s := WrapIndices("Nothing", []int{})
	assert.Equal(t, "Nothing", s)
}

func TestWrapIndices(t *testing.T) {
	s := WrapIndices("abcdefg", []int{1, 2, 5, 6})
	assert.Contains(t, s, "cde")
	assert.Contains(t, s, Red)
	assert.Contains(t, s, Green)
	assert.Contains(t, s, Reset)
}
