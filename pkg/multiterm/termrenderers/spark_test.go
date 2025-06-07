package termrenderers

import (
	"testing"

	"github.com/zix99/rare/pkg/aggregation"
	"github.com/zix99/rare/pkg/aggregation/sorting"
	"github.com/zix99/rare/pkg/multiterm"

	"github.com/stretchr/testify/assert"
)

func TestSimpleSpark(t *testing.T) {
	vt := multiterm.NewVirtualTerm()
	s := NewSpark(vt, 2, 2)

	agg := aggregation.NewTable(" ")
	agg.Sample("a 1")
	agg.Sample("a 2")

	s.WriteTable(agg, sorting.NVNameSorter, sorting.NVNameSorter)
	s.WriteFooter(0, "hello")

	assert.Equal(t, "  First aa Last ", vt.Get(0))
	assert.Equal(t, "1 1     _  1    ", vt.Get(1))
	assert.Equal(t, "2 1     _  1    ", vt.Get(2))
	assert.Equal(t, "hello", vt.Get(3))
	assert.Equal(t, "", vt.Get(4))

	s.Close()
	assert.True(t, vt.IsClosed())
}

func TestOverflowSpark(t *testing.T) {
	vt := multiterm.NewVirtualTerm()
	s := NewSpark(vt, 2, 2)

	agg := aggregation.NewTable(" ")
	agg.Sample("1 a")
	agg.Sample("2 a")
	agg.Sample("2 b")
	agg.Sample("2 b")
	agg.Sample("1 c")

	s.WriteTable(agg, sorting.NVNameSorter, sorting.NVNameSorter)
	s.WriteFooter(0, "hello")

	assert.Equal(t, "  First 12 Last ", vt.Get(0))
	assert.Equal(t, "a 1     ▄▄ 1    ", vt.Get(1))
	assert.Equal(t, "b 0     _█ 2    ", vt.Get(2))
	assert.Equal(t, "(1 more)", vt.Get(3))
	assert.Equal(t, "hello", vt.Get(4))
	assert.Equal(t, "", vt.Get(5))

	s.Close()
	assert.True(t, vt.IsClosed())
}

func TestEmptySpark(t *testing.T) {
	vt := multiterm.NewVirtualTerm()
	s := NewSpark(vt, 2, 2)

	agg := aggregation.NewTable(" ")
	s.WriteTable(agg, sorting.NVNameSorter, sorting.NVNameSorter)

	s.Close()
	assert.True(t, vt.IsClosed())
}
