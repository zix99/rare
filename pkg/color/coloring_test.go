package color

import (
	"fmt"
	"io"
	"strings"
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

func BenchmarkColorReplacerOverlapping(b *testing.B) {
	s := "This is a test"
	groups := []int{4, 7, 5, 6, 8, 9}

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

func TestWrapIndicesInnerGroups(t *testing.T) {
	s := WrapIndices("abcdefg", []int{0, 2, 1, 2, 5, 6})
	assert.Contains(t, s, "cde")
}

func TestWriteColor(t *testing.T) {
	var sb strings.Builder
	Write(&sb, Red, func(w io.StringWriter) {
		w.WriteString("hi")
	})
	assert.Contains(t, sb.String(), "hi")
}

func TestLookupColor(t *testing.T) {
	var c ColorCode
	var ok bool

	c, ok = LookupColorByName("red")
	assert.NotEmpty(t, c)
	assert.True(t, ok)

	c, ok = LookupColorByName("Red")
	assert.NotEmpty(t, c)
	assert.True(t, ok)

	c, ok = LookupColorByName("bla")
	assert.NotEmpty(t, c)
	assert.False(t, ok)
}

func TestHighlightSingleRune(t *testing.T) {
	assert.Equal(t, "\x1b[34;1m\x1b[0m", HighlightSingleRune("", 0, BrightBlue, Underline))
	assert.Equal(t, "\x1b[34;1m\x1b[4m\x1b[36;1ma\x1b[0m\x1b[34;1mbc\x1b[0m", HighlightSingleRune("abc", 0, BrightBlue, Underline+BrightCyan))
	assert.Equal(t, "\x1b[34;1ma\x1b[4m\x1b[36;1mb\x1b[0m\x1b[34;1mc\x1b[0m", HighlightSingleRune("abc", 1, BrightBlue, Underline+BrightCyan))
	assert.Equal(t, "\x1b[34;1mab\x1b[4m\x1b[36;1mc\x1b[0m\x1b[34;1m\x1b[0m", HighlightSingleRune("abc", 2, BrightBlue, Underline+BrightCyan))
	assert.Equal(t, "\x1b[34;1mabc\x1b[0m", HighlightSingleRune("abc", 3, BrightBlue, Underline+BrightCyan))                                   // Out of bounds
	assert.Equal(t, "\x1b[34;1m✤b\x1b[4m\x1b[36;1mc\x1b[0m\x1b[34;1m\x1b[0m", HighlightSingleRune("✤bc", 2, BrightBlue, Underline+BrightCyan)) // Unicode

	Enabled = false
	assert.Equal(t, "test", HighlightSingleRune("test", 1, BrightBlue, BrightCyan))
	Enabled = true
}

func TestStringLength(t *testing.T) {
	assert.Equal(t, 4, StrLen("test"))
	assert.Equal(t, 4, StrLen(Wrap(Red, "test")))
	assert.Equal(t, 12, StrLen(Wrap(Underline, Wrap(Red, "test")+" another")))
	assert.Equal(t, 4, StrLen(Wrap(Red, Wrap(Yellow, "test"))))
	assert.Equal(t, 3, StrLen(Wrap(Red, "ab✤")))

	Enabled = false
	assert.Equal(t, 4, StrLen(Wrap(Red, "test")))
	assert.Equal(t, 3, StrLen("ab✤"))
	Enabled = true
}
