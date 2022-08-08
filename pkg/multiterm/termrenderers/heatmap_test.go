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

	hm.WriteTable(agg)

	assert.Equal(t, 3, vt.LineCount())
	assert.Equal(t, "        █ 0    █ 0    █ 1", vt.Get(0))
	assert.Equal(t, "     test", vt.Get(1))
	assert.Equal(t, "abc  █", vt.Get(2))
	assert.Equal(t, "", vt.Get(3))
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
