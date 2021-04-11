package termrenderers

import (
	"rare/pkg/multiterm"
	"testing"
)

func TestBargraphRendering(t *testing.T) {
	v := multiterm.NewVirtualTerm()
	bg := NewBarGraph(v)
	bg.Stacked = false

	bg.SetKeys("a", "b")
	bg.WriteBar(0, "test", 1, 2)
	bg.WriteBar(1, "tes2", 1, 2)
	bg.WriteFooter(0, "abc")
}

func TestBargraphStackedRendering(t *testing.T) {
	v := multiterm.NewVirtualTerm()
	bg := NewBarGraph(v)
	bg.Stacked = true

	bg.SetKeys("a", "b")
	bg.WriteBar(0, "test", 1, 2)
	bg.WriteBar(1, "tes2", 1, 2)
	bg.WriteFooter(0, "abc")
}

func TestBargraphBadSubkeys(t *testing.T) {
	v := multiterm.NewVirtualTerm()
	bg := NewBarGraph(v)
	bg.Stacked = false

	bg.WriteBar(0, "test", 1, 2)
	bg.WriteBar(1, "tes2", 1, 2)
	bg.WriteBar(2, "tes2")
	bg.WriteFooter(0, "abc")
}
