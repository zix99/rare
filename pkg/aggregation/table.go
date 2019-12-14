package aggregation

import (
	"rare/pkg/stringSplitter"
	"sort"
	"strconv"
)

type TableRow struct {
	cols map[string]int64
	name string
	sum  int64
}

type TableAggregator struct {
	delim  string
	errors uint64
	rows   map[string]*TableRow
	cols   map[string]uint64 // Columns that track usage count (to sort)
}

func NewTable(delim string) *TableAggregator {
	return &TableAggregator{
		delim:  delim,
		errors: 0,
		rows:   make(map[string]*TableRow, 0),
		cols:   make(map[string]uint64, 0),
	}
}

func (s *TableAggregator) ParseErrors() uint64 {
	return s.errors
}

// Samples item like "<column><delim><row><delim><count>"
func (s *TableAggregator) Sample(ele string) {
	splitter := stringSplitter.Splitter{
		S:     ele,
		Delim: s.delim,
	}
	part0 := splitter.Next()
	part1, has1 := splitter.NextOk()
	part2, has2 := splitter.NextOk()
	if has2 {
		inc, err := strconv.ParseInt(part2, 10, 64)
		if err != nil {
			s.errors++
		} else {
			s.SampleItem(part0, part1, inc)
		}
	} else if has1 {
		s.SampleItem(part0, part1, 1)
	} else {
		s.errors++
	}
}

func (s *TableAggregator) SampleItem(colKey, rowKey string, inc int64) {
	s.cols[colKey]++

	row := s.rows[rowKey]
	if row == nil {
		row = &TableRow{
			cols: make(map[string]int64, 0),
			name: rowKey,
		}
		s.rows[rowKey] = row
	}

	row.cols[colKey] += inc
	row.sum += inc
}

func (s *TableAggregator) ColumnCount() int {
	return len(s.cols)
}

func (s *TableAggregator) Columns() []string {
	keys := make([]string, 0, len(s.cols))
	for k := range s.cols {
		keys = append(keys, k)
	}
	return keys
}

func (s *TableAggregator) OrderedColumns() []string {
	keys := s.Columns()

	sort.Slice(keys, func(i, j int) bool {
		c0 := s.cols[keys[i]]
		c1 := s.cols[keys[j]]
		if c0 == c1 {
			return keys[i] < keys[j]
		}
		return c0 > c1
	})

	return keys
}

func (s *TableAggregator) RowCount() int {
	return len(s.rows)
}

func (s *TableAggregator) Rows() []*TableRow {
	rows := make([]*TableRow, 0, len(s.rows))

	for _, v := range s.rows {
		rows = append(rows, v)
	}

	return rows
}

func (s *TableAggregator) OrderedRows() []*TableRow {
	rows := s.Rows()

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].sum == rows[j].sum {
			return rows[i].name < rows[j].name
		}
		return rows[i].sum > rows[j].sum
	})

	return rows
}

func (s *TableAggregator) OrderedRowsByName() []*TableRow {
	rows := s.Rows()

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].name > rows[j].name
	})

	return rows
}

func (s *TableRow) Name() string {
	return s.name
}

func (s *TableRow) Value(colKey string) int64 {
	return s.cols[colKey]
}
