package termrenderers

import (
	"rare/pkg/multiterm"
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

	assert.Equal(t, "        █ a  █ b", v.Get(0))
	assert.Equal(t, "test  █████████████████████████ 1", v.Get(1))
	assert.Equal(t, "      ██████████████████████████████████████████████████ 2", v.Get(2))
	assert.Equal(t, "tes2  ████████████████▊ 1", v.Get(3))
	assert.Equal(t, "      ██████████████████████████████████████████████████ 3", v.Get(4))
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

	assert.Equal(t, "        █ a  █ b", v.Get(0))
	assert.Equal(t, "test  █████████████████████████████████████████████████  3", v.Get(1))
	assert.Equal(t, "tes2  █████████████████████████████████████████████████  3", v.Get(2))
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

	assert.Equal(t, "        █ a", v.Get(0))
	assert.Equal(t, "test  █████████████████████████ 1", v.Get(1))
	assert.Equal(t, "tes2  ████████████████▊ 1", v.Get(2))
	assert.Equal(t, "      ██████████████████████████████████████████████████ 3", v.Get(3))
	assert.Equal(t, "abc", v.Get(4))
}

func TestBargraphUnicode(t *testing.T) {
	v := multiterm.NewVirtualTerm()
	bg := NewBarGraph(v)
	bg.Stacked = true
	bg.maxKeyLength = 0

	bg.SetKeys("✤", "b")
	bg.WriteBar(0, "test", 1, 2)
	bg.WriteBar(1, "tes2", 1, 2)
	bg.WriteBar(2, "✤✥✦a", 1, 2)
	bg.WriteFooter(0, "abc")

	assert.Equal(t, "    █ ✤  █ b", v.Get(0))
	assert.Equal(t, "test  █████████████████████████████████████████████████  3", v.Get(1))
	assert.Equal(t, "tes2  █████████████████████████████████████████████████  3", v.Get(2))
	assert.Equal(t, "✤✥✦a  █████████████████████████████████████████████████  3", v.Get(3))
	assert.Equal(t, "abc", v.Get(4))
}
