package termrenderers

import (
	"rare/pkg/aggregation"
	"rare/pkg/color"
	"rare/pkg/multiterm"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleHeatmap(t *testing.T) {
	vt := multiterm.NewVirtualTerm()
	hm := NewHeatmap(vt, 10, 10)

	agg := aggregation.NewTable(" ")
	agg.Sample("test abc")

	hm.maxRowKeyWidth = 4
	hm.WriteTable(agg)

	assert.Equal(t, 3, vt.LineCount())
	assert.Equal(t, "     - 1    - 1    - 1", vt.Get(0))
	assert.Equal(t, "     test", vt.Get(1))
	assert.Equal(t, "abc  -", vt.Get(2))
	assert.Equal(t, "", vt.Get(3))
}

func TestCompressedHeatmap(t *testing.T) {
	vt := multiterm.NewVirtualTerm()
	hm := NewHeatmap(vt, 2, 2)

	agg := aggregation.NewTable(" ")
	agg.Sample("test abc")
	agg.Sample("test1 abc")
	agg.Sample("test2 abc")
	agg.Sample("test32323 abc")
	agg.Sample("test abc1")
	agg.Sample("test abc2")
	agg.Sample("test abc3")
	agg.Sample("test abc4")

	hm.maxRowKeyWidth = 4
	hm.WriteTable(agg)
	hm.WriteFooter(0, "footer")

	assert.Equal(t, 6, vt.LineCount())
	assert.Equal(t, "     - 0    - 0    9 1", vt.Get(0))
	assert.Equal(t, "     test (2 more)", vt.Get(1))
	assert.Equal(t, "abc  99", vt.Get(2))
	assert.Equal(t, "abc1 9-", vt.Get(3))
	assert.Equal(t, "(3 more)", vt.Get(4))
	assert.Equal(t, "footer", vt.Get(5))
}

func TestHeatmapHeader(t *testing.T) {
	vt := multiterm.NewVirtualTerm()
	hm := NewHeatmap(vt, 1, 10)

	hm.WriteHeader()
	assert.Equal(t, " ", vt.Get(1))

	hm.WriteHeader("abc")
	assert.Equal(t, " abc", vt.Get(1))

	hm.WriteHeader("abc", "efg")
	assert.Equal(t, " abc", vt.Get(1))

	hm.WriteHeader("abc", "fi0", "fi1", "fi2", "fi3", "fi4", "fi5", "fi6", "fi7", "efg")
	assert.Equal(t, " abc....efg", vt.Get(1))

	hm.WriteHeader("a", "b", "c", "d", "e", "f", "g", "h", "i", "j")
	assert.Equal(t, " a..d..g..j", vt.Get(1))

	hm.WriteHeader("a", "b", "c", "d", "e", "f", "g", "h", "i", "jack")
	assert.Equal(t, " a..d..g..jack", vt.Get(1))

	hm.WriteHeader("a", "b", "c", "d", "e", "f", "gar", "h", "i", "jack")
	assert.Equal(t, " a..d..jack", vt.Get(1))

	hm.WriteHeader("a", "b", "c", "d", "e", "f", "ga", "h", "i", "j")
	assert.Equal(t, " a..d.....j", vt.Get(1))

	hm.WriteHeader("aa", "bb", "cc", "dd", "ee", "ff", "gg", "hh", "ii", "jj")
	assert.Equal(t, " aa..ee..jj", vt.Get(1))

	// 2 more
	hm.WriteHeader("a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l")
	assert.Equal(t, " a..d..g..j (2 more)", vt.Get(1))

	// short, by slightly more
	hm.WriteHeader("abc", "d", "e", "f", "g")
	assert.Equal(t, " abc..", vt.Get(1))

	hm.WriteHeader("abc", "d", "e", "f")
	assert.Equal(t, " abc.", vt.Get(1))
}

func TestUnderlineHeaderChar(t *testing.T) {
	color.Enabled = true
	assert.Equal(t, "\x1b[34;1m\x1b[0m", underlineHeaderChar("", 0))
	assert.Equal(t, "\x1b[34;1m\x1b[0m\x1b[4m\x1b[36;1ma\x1b[0m\x1b[34;1mbc\x1b[0m", underlineHeaderChar("abc", 0))
	assert.Equal(t, "\x1b[34;1ma\x1b[0m\x1b[4m\x1b[36;1mb\x1b[0m\x1b[34;1mc\x1b[0m", underlineHeaderChar("abc", 1))
	assert.Equal(t, "\x1b[34;1mab\x1b[0m\x1b[4m\x1b[36;1mc\x1b[0m\x1b[34;1m\x1b[0m", underlineHeaderChar("abc", 2))
	assert.Equal(t, "\x1b[34;1mabc\x1b[0m", underlineHeaderChar("abc", 3))
	color.Enabled = false
}
