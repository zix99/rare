package termrenderers

import (
	"rare/pkg/multiterm"
	"rare/pkg/multiterm/termscaler"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBargraphRendering(t *testing.T) {
	v := multiterm.NewVirtualTerm()
	bg := NewBarGraph(v)
	bg.Stacked = false

	bg.SetKeys("a", "b")
	bg.WriteBar(0, "test", 1, 2)
	bg.WriteBar(1, "tes2", 1, 3)
	bg.WriteFooter(0, "abc")

	assert.Equal(t, "        0 a  1 b", v.Get(0))
	assert.Equal(t, "test  ████████████████▊ 1", v.Get(1))
	assert.Equal(t, "      █████████████████████████████████▍ 2", v.Get(2))
	assert.Equal(t, "tes2  ████████████████▊ 1", v.Get(3))
	assert.Equal(t, "      ██████████████████████████████████████████████████ 3", v.Get(4))
	assert.Equal(t, "abc", v.Get(5))
}

func TestBargraphRenderingLog10(t *testing.T) {
	v := multiterm.NewVirtualTerm()
	bg := NewBarGraph(v)
	bg.Stacked = false
	bg.Scaler = termscaler.ScalerLog10

	bg.SetKeys("a", "b")
	bg.WriteBar(0, "test", 10, 100)
	bg.WriteBar(1, "tes2", 50, 500)
	bg.WriteFooter(0, "abc")

	assert.Equal(t, "        0 a  1 b", v.Get(0))
	assert.Equal(t, "test  ████████████████▊ 10", v.Get(1))
	assert.Equal(t, "      █████████████████████████████████▍ 100", v.Get(2))
	assert.Equal(t, "tes2  ████████████████████████████▎ 50", v.Get(3))
	assert.Equal(t, "      █████████████████████████████████████████████ 500", v.Get(4))
	assert.Equal(t, "abc", v.Get(5))
}

func TestBargraphStackedRendering(t *testing.T) {
	v := multiterm.NewVirtualTerm()
	bg := NewBarGraph(v)
	bg.Stacked = true

	bg.SetKeys("a", "b")
	bg.WriteBar(0, "test", 1, 2)
	bg.WriteBar(1, "tes2", 1, 2)
	bg.WriteFooter(0, "abc")

	assert.Equal(t, "        0 a  1 b", v.Get(0))
	assert.Equal(t, "test  0000000000000000111111111111111111111111111111111  3", v.Get(1))
	assert.Equal(t, "tes2  0000000000000000111111111111111111111111111111111  3", v.Get(2))
	assert.Equal(t, "abc", v.Get(3))
}

func TestBargraphBadSubkeys(t *testing.T) {
	v := multiterm.NewVirtualTerm()
	bg := NewBarGraph(v)
	bg.Stacked = false

	bg.SetKeys("a")
	bg.WriteBar(0, "test", 1, 2)
	bg.WriteBar(1, "tes2", 1, 3)
	bg.WriteBar(2, "tes2")
	bg.WriteFooter(0, "abc")

	assert.Equal(t, "        0 a", v.Get(0))
	assert.Equal(t, "test  ████████████████▊ 1", v.Get(1))
	assert.Equal(t, "tes2  ████████████████▊ 1", v.Get(2))
	assert.Equal(t, "      ██████████████████████████████████████████████████ 3", v.Get(3))
	assert.Equal(t, "abc", v.Get(4))
}

func TestBargraphUnicodeStacked(t *testing.T) {
	v := multiterm.NewVirtualTerm()
	bg := NewBarGraph(v)
	bg.Stacked = true
	bg.maxKeyLength = 0

	bg.SetKeys("✤", "b")
	bg.WriteBar(0, "test", 1, 2)
	bg.WriteBar(1, "tes2", 1, 2)
	bg.WriteBar(2, "✤✥✦a", 1, 2)
	bg.WriteFooter(0, "abc")

	assert.Equal(t, "    0 ✤  1 b", v.Get(0))
	assert.Equal(t, "test  0000000000000000111111111111111111111111111111111  3", v.Get(1))
	assert.Equal(t, "tes2  0000000000000000111111111111111111111111111111111  3", v.Get(2))
	assert.Equal(t, "✤✥✦a  0000000000000000111111111111111111111111111111111  3", v.Get(3))
	assert.Equal(t, "abc", v.Get(4))
}

func TestMaxSumUtils(t *testing.T) {
	assert.Equal(t, int64(5), maxi64(1, 2, 3, 5, 4))
	assert.Equal(t, int64(15), sumi64(1, 2, 3, 5, 4))
}
