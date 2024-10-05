package aggregation

import (
	"math"
	"rare/pkg/aggregation/sorting"
	"rare/pkg/stringSplitter"
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
	cols   map[string]int64 // Columns that track totals
}

func NewTable(delim string) *TableAggregator {
	return &TableAggregator{
		delim:  delim,
		errors: 0,
		rows:   make(map[string]*TableRow),
		cols:   make(map[string]int64),
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
		s.SampleItem(part0, "", 1)
	}
}

func (s *TableAggregator) SampleItem(colKey, rowKey string, inc int64) {
	s.cols[colKey] += inc

	row := s.rows[rowKey]
	if row == nil {
		row = &TableRow{
			cols: make(map[string]int64),
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

func (s *TableAggregator) OrderedColumns(sorter sorting.NameValueSorter) []string {
	keys := s.Columns()
	sorting.SortBy(keys, sorter, func(name string) sorting.NameValuePair {
		return sorting.NameValuePair{
			Name:  name,
			Value: s.cols[name],
		}
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

func (s *TableAggregator) OrderedRows(sorter sorting.NameValueSorter) []*TableRow {
	rows := s.Rows()
	sorting.SortBy(rows, sorter, func(obj *TableRow) sorting.NameValuePair {
		return sorting.NameValuePair{
			Name:  obj.name,
			Value: obj.sum,
		}
	})
	return rows
}

func (s *TableAggregator) ComputeMin() (ret int64) {
	ret = math.MaxInt64
	for _, r := range s.rows {
		for colKey := range s.cols {
			if val := r.cols[colKey]; val < ret {
				ret = val
			}
		}
	}
	if ret == math.MaxInt64 {
		return 0
	}
	return
}

func (s *TableAggregator) ComputeMax() (ret int64) {
	ret = math.MinInt64
	for _, r := range s.rows {
		for colKey := range s.cols {
			if val := r.cols[colKey]; val > ret {
				ret = val
			}
		}
	}
	if ret == math.MinInt64 {
		return 0
	}
	return
}

// ColTotals returns column oriented totals (Do not change!)
func (s *TableAggregator) ColTotal(k string) int64 {
	return s.cols[k]
}

// Sum all data
func (s *TableAggregator) Sum() (ret int64) {
	for _, v := range s.cols {
		ret += v
	}
	return
}

// Trim data. Returns number of fields trimmed
func (s *TableAggregator) Trim(predicate func(col, row string, val int64) bool) int {
	trimmed := 0

	// TODO: Ability to delete data from the table based on predicate

	return trimmed
}

func (s *TableRow) Name() string {
	return s.name
}

func (s *TableRow) Value(colKey string) int64 {
	return s.cols[colKey]
}

// Sum is the total sum of all values in the row
func (s *TableRow) Sum() int64 {
	return s.sum
}
