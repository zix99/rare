package csv

import (
	"strconv"

	"github.com/zix99/rare/pkg/aggregation"
	"github.com/zix99/rare/pkg/aggregation/sorting"
)

func WriteTable(w CSV, agg *aggregation.TableAggregator) error {
	cols := agg.OrderedColumns(sorting.NVNameSorter)
	w.Write(append([]string{""}, cols...))

	rowBuf := make([]string, len(cols)+1)
	for _, row := range agg.OrderedRows(sorting.NVNameSorter) {
		rowBuf[0] = row.Name()
		for i, colName := range cols {
			rowBuf[i+1] = strconv.FormatInt(row.Value(colName), 10)
		}
		if err := w.Write(rowBuf); err != nil {
			return err
		}
	}

	return nil
}

func WriteAccumulator(w CSV, aggr *aggregation.AccumulatingGroup) error {
	{
		header := make([]string, 0, aggr.ColCount())
		header = append(header, aggr.GroupCols()...)
		header = append(header, aggr.DataCols()...)
		w.Write(header)
	}

	row := make([]string, aggr.ColCount())
	for _, group := range aggr.Groups(sorting.ByName) {
		copy(row, group.Parts())
		copy(row[aggr.GroupColCount():], aggr.DataNoCopy(group))
		if err := w.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func WriteCounter(w CSV, aggr *aggregation.MatchCounter) error {
	w.Write([]string{"group", "value"})

	for _, group := range aggr.ItemsSortedBy(aggr.GroupCount(), sorting.NVValueSorter) {
		if err := w.WriteRow(group.Name, group.Item.Count()); err != nil {
			return err
		}
	}

	return nil
}

func WriteSubCounter(w CSV, aggr *aggregation.SubKeyCounter) error {
	header := append([]string{"group"}, aggr.SubKeys()...)
	w.Write(header)

	row := make([]string, len(header))
	for _, item := range aggr.ItemsSorted(sorting.NVNameSorter) {
		row[0] = item.Name
		for i, val := range item.Item.Items() {
			row[i+1] = strconv.FormatInt(val, 10)
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}
	return nil
}
