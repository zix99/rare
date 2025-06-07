package termrenderers

import (
	"github.com/zix99/rare/pkg/aggregation"
	"github.com/zix99/rare/pkg/aggregation/sorting"
	"github.com/zix99/rare/pkg/color"
	"github.com/zix99/rare/pkg/multiterm"
	"github.com/zix99/rare/pkg/multiterm/termformat"
)

type DataTable struct {
	table            *TableWriter
	numRows, numCols int
	ShowRowTotals    bool
	ShowColTotals    bool

	formatter   termformat.Formatter
	needsMinMax bool
}

func NewDataTable(term multiterm.MultilineTerm, numCols, numRows int) *DataTable {
	return &DataTable{
		table:     NewTable(term, numCols+2, numRows+2),
		numRows:   numRows,
		numCols:   numCols,
		formatter: termformat.Default,
	}
}

func (s *DataTable) SetFormatter(f termformat.Formatter) {
	s.formatter = f
	s.needsMinMax = true
}

func (s *DataTable) WriteTable(counter *aggregation.TableAggregator, rowSorter, colSorter sorting.NameValueSorter) {
	cols := counter.OrderedColumns(colSorter)
	cols = minColSlice(s.numCols, cols) // Cap columns

	var min, max int64
	if s.needsMinMax {
		min, max = counter.ComputeMinMax()
	}

	// Write header row
	{
		colNames := make([]string, len(cols)+2)
		for i, name := range cols {
			colNames[i+1] = color.Wrap(color.Underline+color.BrightBlue, name)
		}
		if s.ShowRowTotals {
			colNames[len(cols)+1] = color.Wrap(color.Underline+color.BrightBlack, "Total")
		}
		s.table.WriteRow(0, colNames...)
	}

	// Write each row
	rows := counter.OrderedRows(rowSorter)

	line := 1
	for i := 0; i < len(rows) && i < s.numRows; i++ {
		row := rows[i]
		rowVals := make([]string, len(cols)+2)
		rowVals[0] = color.Wrap(color.Yellow, row.Name())
		for idx, colName := range cols {
			rowVals[idx+1] = s.formatter(row.Value(colName), min, max)
		}
		if s.ShowRowTotals {
			rowVals[len(rowVals)-1] = color.Wrap(color.BrightBlack, s.formatter(row.Sum(), min, max))
		}
		s.table.WriteRow(line, rowVals...)
		line++
	}

	// Write totals
	if s.ShowColTotals {
		rowVals := make([]string, len(cols)+2)
		rowVals[0] = color.Wrap(color.BrightBlack+color.Underline, "Total")
		for idx, colName := range cols {
			rowVals[idx+1] = color.Wrap(color.BrightBlack, s.formatter(counter.ColTotal(colName), min, max))
		}

		if s.ShowRowTotals { // super total
			sum := counter.Sum()
			rowVals[len(rowVals)-1] = color.Wrap(color.BrightWhite, s.formatter(sum, min, max))
		}

		s.table.WriteRow(line, rowVals...)
	}
}

func (s *DataTable) Close() {
	s.table.Close()
}

func (s *DataTable) WriteFooter(idx int, line string) {
	s.table.WriteFooter(idx, line)
}

func minColSlice(count int, cols []string) []string {
	if len(cols) < count {
		return cols
	}
	return cols[:count]
}
