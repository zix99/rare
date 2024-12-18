package aggregation

import (
	"fmt"
	"rare/pkg/aggregation/sorting"
	"strconv"
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

	assert.Equal(t, []string{"b", "a"}, table.OrderedColumns(sorting.NVValueSorter))

	rows := table.OrderedRows(sorting.NVValueSorter)
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
	assert.Equal(t, []string{"b", "a"}, table.OrderedColumns(sorting.NVValueSorter))
	assert.Equal(t, uint64(1), table.ParseErrors())

	// Col totals
	assert.Equal(t, int64(2), table.ColTotal("a"))
	assert.Equal(t, int64(3), table.ColTotal("b"))
	assert.Equal(t, int64(0), table.ColTotal("c"))

	// Totals
	assert.Equal(t, int64(5), table.Sum())

	// Minmax
	min, max := table.ComputeMinMax()
	assert.Equal(t, int64(0), min)
	assert.Equal(t, int64(3), max)
}

func TestTableMultiIncrement(t *testing.T) {
	table := NewTable(" ")
	table.Sample("a b 1")
	table.Sample("b c 3")
	table.Sample("b c 3")
	table.Sample("b c -1")

	// Row names and col vals
	rows := table.OrderedRows(sorting.NVNameSorter)
	assert.Equal(t, "b", rows[0].Name())
	assert.Equal(t, int64(1), rows[0].Value("a"))
	assert.Equal(t, "c", rows[1].Name())
	assert.Equal(t, int64(5), rows[1].Value("b"))

	// Column names
	assert.Equal(t, []string{"a", "b"}, table.OrderedColumns(sorting.NVNameSorter))
	assert.Equal(t, []string{"b", "a"}, table.OrderedColumns(sorting.NVValueSorter))

	// Totals
	assert.Equal(t, int64(5), table.ColTotal("b"))
	assert.Equal(t, int64(1), table.ColTotal("a"))
	assert.Equal(t, int64(6), table.Sum())

	// Minmax
	min, max := table.ComputeMinMax()
	assert.Equal(t, int64(0), min)
	assert.Equal(t, int64(5), max)
}

func TestEmptyTableMinMax(t *testing.T) {
	table := NewTable(" ")
	min, max := table.ComputeMinMax()
	assert.Equal(t, int64(0), min)
	assert.Equal(t, int64(0), max)
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

func TestTrimData(t *testing.T) {
	table := NewTable(" ")
	for i := 0; i < 10; i++ {
		table.Sample(fmt.Sprintf("%d a", i))
		table.Sample(fmt.Sprintf("%d b", i))
	}

	assert.Len(t, table.Columns(), 10)

	trimmed := table.Trim(func(col, row string, val int64) bool {
		if row == "b" {
			return true
		}
		cVal, _ := strconv.Atoi(col)
		return cVal < 5
	})

	assert.ElementsMatch(t, []string{"5", "6", "7", "8", "9"}, table.Columns())
	assert.Equal(t, 15, trimmed)
	assert.Len(t, table.Rows(), 1)
	assert.Len(t, table.Rows()[0].cols, 5)
}

// BenchmarkMinMax-4   	 1020728	      1234 ns/op	       0 B/op	       0 allocs/op
func BenchmarkMinMax(b *testing.B) {
	table := NewTable(" ")
	for i := 0; i < 10; i++ {
		table.Sample(fmt.Sprintf("%d a", i))
		table.Sample(fmt.Sprintf("%d b", i))
	}

	for i := 0; i < b.N; i++ {
		table.ComputeMinMax()
	}
}
