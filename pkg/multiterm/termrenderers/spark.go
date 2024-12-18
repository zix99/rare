package termrenderers

import (
	"rare/pkg/aggregation"
	"rare/pkg/aggregation/sorting"
	"rare/pkg/color"
	"rare/pkg/humanize"
	"rare/pkg/multiterm"
	"rare/pkg/multiterm/termscaler"
	"rare/pkg/multiterm/termunicode"
	"strings"
)

type Spark struct {
	rowCount, colCount int
	footerOffset       int
	Scaler             termscaler.Scaler
	table              *TableWriter
}

func NewSpark(term multiterm.MultilineTerm, rows, cols int) *Spark {
	return &Spark{
		rowCount: rows,
		colCount: cols,
		Scaler:   termscaler.ScalerLinear,
		table:    NewTable(term, 4, rows+1),
	}
}

func (s *Spark) WriteTable(agg *aggregation.TableAggregator, rowSorter, colSorter sorting.NameValueSorter) {
	minVal, maxVal := agg.ComputeMinMax()

	colNames := agg.OrderedColumns(colSorter)
	if len(colNames) > s.colCount {
		colNames = colNames[len(colNames)-s.colCount:]
	}

	// reused buffer
	var sb strings.Builder
	sb.Grow(len(colNames))

	// Write header
	{
		dots := len(colNames) - len(colNames[0]) - len(colNames[len(colNames)-1])
		if dots < 0 {
			dots = 0
		}
		sb.WriteString(colNames[0])
		writeRepeat(&sb, '.', dots)
		sb.WriteString(colNames[len(colNames)-1])

		s.table.WriteRow(0, "", color.Wrap(color.Underline, "First"), sb.String(), color.Wrap(color.Underline, "Last"))
		sb.Reset()
	}

	// Each row...
	rows := agg.OrderedRows(rowSorter)
	rowCount := mini(len(rows), s.rowCount)
	for i := 0; i < rowCount; i++ {
		row := rows[i]

		for j := 0; j < len(colNames); j++ {
			termunicode.SparkWrite(&sb, s.Scaler.Scale(row.Value(colNames[j]), minVal, maxVal))
		}

		vFirst := humanize.Hi(row.Value(colNames[0]))
		vLast := humanize.Hi(row.Value(colNames[len(colNames)-1]))
		s.table.WriteRow(i+1, color.Wrap(color.Yellow, row.Name()), color.Wrap(color.BrightBlack, vFirst), sb.String(), color.Wrap(color.BrightBlack, vLast))

		sb.Reset()
	}

	// If more rows than can display, write how many were missed
	if len(rows) > rowCount {
		s.table.WriteFooter(0, color.Wrapf(color.BrightBlack, "(%d more)", len(rows)-rowCount))
		s.footerOffset = 1
	} else {
		s.footerOffset = 0
	}
}

func (s *Spark) Close() {
	s.table.Close()
}

func (s *Spark) WriteFooter(idx int, line string) {
	s.table.WriteFooter(s.footerOffset+idx, line)
}
