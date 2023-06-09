package csv

import (
	"rare/pkg/aggregation"
	"rare/pkg/aggregation/sorting"
	"strconv"
)

func WriteTable(w CSV, agg *aggregation.TableAggregator) error {
	cols := agg.OrderedColumns(sorting.NVNameSorter)
	w.Write(append([]string{""}, cols...))
	for _, row := range agg.OrderedRows(sorting.NVNameSorter) {
		arr := make([]string, len(cols)+1)
		arr[0] = row.Name()
		for i, colName := range cols {
			arr[i+1] = strconv.FormatInt(row.Value(colName), 10)
		}
		if err := w.Write(arr); err != nil {
			return err
		}
	}

	return nil
}

func WriteAccumulator(w CSV, aggr *aggregation.AccumulatingGroup) error {
	{
		header := make([]string, 0)
		header = append(header, aggr.GroupCols()...)
		header = append(header, aggr.DataCols()...)
		w.Write(header)
	}

	for _, group := range aggr.Groups(sorting.ByName) {
		row := make([]string, 0)
		row = append(row, group.Parts()...)
		row = append(row, aggr.Data(group)...)
		if err := w.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func WriteCounter(w CSV, aggr *aggregation.MatchCounter) error {
	w.WriteRow("group", "value")

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

	for _, item := range aggr.ItemsSorted(sorting.NVNameSorter) {
		row := []any{item.Name}
		for _, val := range item.Item.Items() {
			row = append(row, val)
		}
		if err := w.WriteRow(row...); err != nil {
			return err
		}
	}
	return nil
}
