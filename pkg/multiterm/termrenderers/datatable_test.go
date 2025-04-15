package termrenderers

import (
	"rare/pkg/aggregation"
	"rare/pkg/aggregation/sorting"
	"rare/pkg/multiterm"
	"rare/pkg/multiterm/termformat"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleDataTable(t *testing.T) {
	vt := multiterm.NewVirtualTerm()
	s := NewDataTable(vt, 2, 2)

	agg := aggregation.NewTable(" ")
	agg.Sample("a 1")
	agg.Sample("a 2")

	s.WriteTable(agg, sorting.NVNameSorter, sorting.NVNameSorter)
	s.WriteFooter(0, "howdy")

	assert.Equal(t, "  a  ", vt.Get(0))
	assert.Equal(t, "1 1  ", vt.Get(1))
	assert.Equal(t, "2 1  ", vt.Get(2))
	assert.Equal(t, "howdy", vt.Get(3))
	assert.Equal(t, "", vt.Get(4))

	s.Close()
	assert.True(t, vt.IsClosed())
}

func TestDataTableFormatter(t *testing.T) {
	vt := multiterm.NewVirtualTerm()
	s := NewDataTable(vt, 2, 2)
	s.SetFormatter(termformat.MustFromExpression("{multi {0} 10}"))

	agg := aggregation.NewTable(" ")
	agg.Sample("a 1")
	agg.Sample("a 2")

	s.WriteTable(agg, sorting.NVNameSorter, sorting.NVNameSorter)
	s.WriteFooter(0, "howdy")

	assert.Equal(t, "  a   ", vt.Get(0))
	assert.Equal(t, "1 10  ", vt.Get(1))
	assert.Equal(t, "2 10  ", vt.Get(2))
	assert.Equal(t, "howdy", vt.Get(3))
	assert.Equal(t, "", vt.Get(4))

	s.Close()
	assert.True(t, vt.IsClosed())
}

func TestComplexDataTable(t *testing.T) {
	vt := multiterm.NewVirtualTerm()
	s := NewDataTable(vt, 3, 3)

	agg := aggregation.NewTable(" ")
	agg.Sample("a 1")
	agg.Sample("a 2")
	agg.Sample("b z")
	agg.Sample("b z")
	agg.Sample("c e")
	agg.Sample("c e")
	agg.Sample("c e")

	s.WriteTable(agg, sorting.NVNameSorter, sorting.NVNameSorter)
	s.WriteFooter(0, "howdy")

	assert.Equal(t, "  a b c  ", vt.Get(0))
	assert.Equal(t, "1 1 0 0  ", vt.Get(1))
	assert.Equal(t, "2 1 0 0  ", vt.Get(2))
	assert.Equal(t, "e 0 0 3  ", vt.Get(3))
	assert.Equal(t, "howdy", vt.Get(4))

	s.Close()
	assert.True(t, vt.IsClosed())
}

func TestDataTableRowTotals(t *testing.T) {
	vt := multiterm.NewVirtualTerm()
	s := NewDataTable(vt, 3, 3)

	agg := aggregation.NewTable(" ")
	agg.Sample("a 1")
	agg.Sample("a 2")
	agg.Sample("b z")
	agg.Sample("b z")
	agg.Sample("c e")
	agg.Sample("c e")
	agg.Sample("c e")

	s.ShowColTotals = false
	s.ShowRowTotals = true
	s.WriteTable(agg, sorting.NVNameSorter, sorting.NVNameSorter)
	s.WriteFooter(0, "howdy")
	assert.Equal(t, "  a b c Total ", vt.Get(0))
	assert.Equal(t, "1 1 0 0 1     ", vt.Get(1))
	assert.Equal(t, "2 1 0 0 1     ", vt.Get(2))
	assert.Equal(t, "e 0 0 3 3     ", vt.Get(3))
	assert.Equal(t, "howdy", vt.Get(4))

	s.Close()
	assert.True(t, vt.IsClosed())
}

func TestDataTableColTotals(t *testing.T) {
	vt := multiterm.NewVirtualTerm()
	s := NewDataTable(vt, 3, 3)

	agg := aggregation.NewTable(" ")
	agg.Sample("a 1")
	agg.Sample("a 2")
	agg.Sample("b z")
	agg.Sample("b z")
	agg.Sample("c e")
	agg.Sample("c e")
	agg.Sample("c e")

	s.ShowColTotals = true
	s.ShowRowTotals = false
	s.WriteTable(agg, sorting.NVNameSorter, sorting.NVNameSorter)
	s.WriteFooter(0, "howdy")
	assert.Equal(t, "      a b c  ", vt.Get(0))
	assert.Equal(t, "1     1 0 0  ", vt.Get(1))
	assert.Equal(t, "2     1 0 0  ", vt.Get(2))
	assert.Equal(t, "e     0 0 3  ", vt.Get(3))
	assert.Equal(t, "Total 2 2 3  ", vt.Get(4))
	assert.Equal(t, "howdy", vt.Get(5))

	s.Close()
	assert.True(t, vt.IsClosed())
}

func TestDataTableRowAndColTotals(t *testing.T) {
	vt := multiterm.NewVirtualTerm()
	s := NewDataTable(vt, 3, 3)

	agg := aggregation.NewTable(" ")
	agg.Sample("a 1")
	agg.Sample("a 2")
	agg.Sample("b z")
	agg.Sample("b z")
	agg.Sample("c e")
	agg.Sample("c e")
	agg.Sample("c e")

	s.ShowColTotals = true
	s.ShowRowTotals = true
	s.WriteTable(agg, sorting.NVNameSorter, sorting.NVNameSorter)
	s.WriteFooter(0, "howdy")
	assert.Equal(t, "      a b c Total ", vt.Get(0))
	assert.Equal(t, "1     1 0 0 1     ", vt.Get(1))
	assert.Equal(t, "2     1 0 0 1     ", vt.Get(2))
	assert.Equal(t, "e     0 0 3 3     ", vt.Get(3))
	assert.Equal(t, "Total 2 2 3 7     ", vt.Get(4))
	assert.Equal(t, "howdy", vt.Get(5))

	s.Close()
	assert.True(t, vt.IsClosed())
}

func TestOverflowDataTable(t *testing.T) {
	vt := multiterm.NewVirtualTerm()
	s := NewDataTable(vt, 2, 2)

	agg := aggregation.NewTable(" ")
	agg.Sample("a 1")
	agg.Sample("a 2")
	agg.Sample("b z")
	agg.Sample("b z")
	agg.Sample("c e")
	agg.Sample("c e")
	agg.Sample("c e")

	s.WriteTable(agg, sorting.NVNameSorter, sorting.NVNameSorter)
	s.WriteFooter(0, "howdy")

	assert.Equal(t, "  a b  ", vt.Get(0))
	assert.Equal(t, "1 1 0  ", vt.Get(1))
	assert.Equal(t, "2 1 0  ", vt.Get(2))
	assert.Equal(t, "howdy", vt.Get(3))

	s.Close()
	assert.True(t, vt.IsClosed())
}

func TestEmptyDataTable(t *testing.T) {
	vt := multiterm.NewVirtualTerm()
	s := NewDataTable(vt, 2, 2)
	agg := aggregation.NewTable(" ")

	s.WriteTable(agg, sorting.NVNameSorter, sorting.NVNameSorter)
	assert.Equal(t, "  ", vt.Get(0))

	s.ShowColTotals = true
	s.ShowRowTotals = true
	s.WriteTable(agg, sorting.NVNameSorter, sorting.NVNameSorter)
	assert.Equal(t, "      Total ", vt.Get(0))

	s.Close()
	assert.True(t, vt.IsClosed())
}
