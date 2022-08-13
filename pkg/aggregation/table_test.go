package aggregation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleTable(t *testing.T) {
	table := NewTable(" ")
	table.Sample("b c")
	table.Sample("b c")
	table.Sample("a b")
	table.Sample("a c")
	table.Sample("b c")
	table.Sample("b b q") // invalid

	assert.Equal(t, []string{"b", "a"}, table.OrderedColumns())

	rows := table.OrderedRows()
	assert.Equal(t, 2, len(rows))
	assert.Equal(t, "c", rows[0].Name())
	assert.Equal(t, "b", rows[1].Name())
	assert.Equal(t, int64(4), rows[0].Sum())
	assert.Equal(t, int64(1), rows[1].Sum())
	assert.Equal(t, int64(3), rows[0].Value("b"))
	assert.Equal(t, int64(1), rows[0].Value("a"))

	assert.Equal(t, 2, table.RowCount())
	assert.Equal(t, 2, table.ColumnCount())

	assert.Contains(t, table.Columns(), "a")
	assert.Contains(t, table.Columns(), "b")
	assert.Equal(t, []string{"b", "a"}, table.OrderedColumns())
	assert.Equal(t, uint64(1), table.ParseErrors())

	// Col totals
	assert.Equal(t, int64(2), table.ColTotal("a"))
	assert.Equal(t, int64(3), table.ColTotal("b"))
	assert.Equal(t, int64(0), table.ColTotal("c"))

	// Totals
	assert.Equal(t, int64(5), table.Sum())

	// Minmax
	assert.Equal(t, int64(0), table.ComputeMin())
	assert.Equal(t, int64(3), table.ComputeMax())
}

func TestTableMultiIncrement(t *testing.T) {
	table := NewTable(" ")
	table.Sample("a b 1")
	table.Sample("b c 3")
	table.Sample("b c 3")
	table.Sample("b c -1")

	// Row names and col vals
	rows := table.OrderedRowsByName()
	assert.Equal(t, "b", rows[0].Name())
	assert.Equal(t, int64(1), rows[0].Value("a"))
	assert.Equal(t, "c", rows[1].Name())
	assert.Equal(t, int64(5), rows[1].Value("b"))

	// Column names
	assert.Equal(t, []string{"a", "b"}, table.OrderedColumnsByName())
	assert.Equal(t, []string{"b", "a"}, table.OrderedColumns())

	// Totals
	assert.Equal(t, int64(5), table.ColTotal("b"))
	assert.Equal(t, int64(1), table.ColTotal("a"))
	assert.Equal(t, int64(6), table.Sum())

	// Minmax
	assert.Equal(t, int64(0), table.ComputeMin())
	assert.Equal(t, int64(5), table.ComputeMax())
}

func TestSingleRowTable(t *testing.T) {
	table := NewTable(" ")
	table.Sample("a")
	table.Sample("b")
	table.Sample("a")

	rows := table.Rows()
	assert.Len(t, rows, 1)
	assert.Empty(t, rows[0].Name())

	assert.Len(t, table.Columns(), 2)

	assert.Equal(t, int64(2), rows[0].Value("a"))
	assert.Equal(t, int64(1), rows[0].Value("b"))
}
