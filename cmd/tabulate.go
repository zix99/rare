package cmd

import (
	"fmt"
	"rare/cmd/helpers"
	"rare/pkg/aggregation"
	"rare/pkg/color"
	"rare/pkg/csv"
	"rare/pkg/expressions"
	"rare/pkg/multiterm/termrenderers"

	"github.com/urfave/cli/v2"
)

func minColSlice(count int, cols []string) []string {
	if len(cols) < count {
		return cols
	}
	return cols[:count]
}

func tabulateFunction(c *cli.Context) error {
	var (
		delim     = c.String("delim")
		numRows   = c.Int("num")
		numCols   = c.Int("cols")
		rowtotals = c.Bool("rowtotal") || c.Bool("x")
		coltotals = c.Bool("coltotal") || c.Bool("x")
		sortRows  = c.String("sort-rows")
		sortCols  = c.String("sort-cols")
		formatExp = c.String(helpers.FormatFlag.Name)
	)

	counter := aggregation.NewTable(delim)
	vt := helpers.BuildVTermFromArguments(c)
	writer := termrenderers.NewTable(vt, numCols+2, numRows+2)

	batcher := helpers.BuildBatcherFromArguments(c)
	ext := helpers.BuildExtractorFromArguments(c, batcher)
	rowSorter := helpers.BuildSorterOrFail(sortRows)
	colSorter := helpers.BuildSorterOrFail(sortCols)
	formatter := helpers.BuildFormatterOrFail(formatExp)

	var min, max int64
	needsMinMax := (formatExp != "")

	helpers.RunAggregationLoop(ext, counter, func() {
		cols := counter.OrderedColumns(colSorter)
		cols = minColSlice(numCols, cols) // Cap columns

		if needsMinMax {
			min, max = counter.ComputeMinMax()
		}

		// Write header row
		{
			colNames := make([]string, len(cols)+2)
			for i, name := range cols {
				colNames[i+1] = color.Wrap(color.Underline+color.BrightBlue, name)
			}
			if rowtotals {
				colNames[len(cols)+1] = color.Wrap(color.Underline+color.BrightBlack, "Total")
			}
			writer.WriteRow(0, colNames...)
		}

		// Write each row
		rows := counter.OrderedRows(rowSorter)

		line := 1
		for i := 0; i < len(rows) && i < numRows; i++ {
			row := rows[i]
			rowVals := make([]string, len(cols)+2)
			rowVals[0] = color.Wrap(color.Yellow, row.Name())
			for idx, colName := range cols {
				rowVals[idx+1] = formatter(row.Value(colName), min, max)
			}
			if rowtotals {
				rowVals[len(rowVals)-1] = color.Wrap(color.BrightBlack, formatter(row.Sum(), min, max))
			}
			writer.WriteRow(line, rowVals...)
			line++
		}

		// Write totals
		if coltotals {
			rowVals := make([]string, len(cols)+2)
			rowVals[0] = color.Wrap(color.BrightBlack+color.Underline, "Total")
			for idx, colName := range cols {
				rowVals[idx+1] = color.Wrap(color.BrightBlack, formatter(counter.ColTotal(colName), min, max))
			}

			if rowtotals { // super total
				sum := counter.Sum()
				rowVals[len(rowVals)-1] = color.Wrap(color.BrightWhite, formatter(sum, min, max))
			}

			writer.WriteRow(line, rowVals...)
		}

		writer.WriteFooter(0, helpers.FWriteExtractorSummary(ext, counter.ParseErrors(),
			fmt.Sprintf("(R: %v; C: %v)", color.Wrapi(color.Yellow, counter.RowCount()), color.Wrapi(color.BrightBlue, counter.ColumnCount()))))
		writer.WriteFooter(1, batcher.StatusString())
	})

	writer.Close()

	if err := helpers.TryWriteCSV(c, counter, csv.WriteTable); err != nil {
		return err
	}

	return helpers.DetermineErrorState(batcher, ext, counter)
}

func tabulateCommand() *cli.Command {
	return helpers.AdaptCommandForExtractor(cli.Command{
		Name:    "tabulate",
		Aliases: []string{"table", "t"},
		Usage:   "Create a 2D summarizing table of extracted data",
		Description: `Summarizes the extracted data as a 2D data table.
		The expression key data format is {$ a b [c]}, where a is the column key,
		b is the rowkey, and optionally c is the increment value (Default: 1)`,
		Action:   tabulateFunction,
		Category: cmdCatVisualize,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "delim",
				Usage: "Character to tabulate on. Use {$} helper by default",
				Value: expressions.ArraySeparatorString,
			},
			&cli.IntFlag{
				Name:    "num",
				Aliases: []string{"rows", "n"},
				Usage:   "Number of elements to display",
				Value:   20,
			},
			&cli.IntFlag{
				Name:  "cols",
				Usage: "Number of columns to display",
				Value: 10,
			},
			&cli.BoolFlag{
				Name:  "rowtotal",
				Usage: "Show row totals",
			},
			&cli.BoolFlag{
				Name:  "coltotal",
				Usage: "Show column totals",
			},
			&cli.BoolFlag{
				Name:    "extra",
				Aliases: []string{"x"},
				Usage:   "Display row and column totals",
			},
			&cli.StringFlag{
				Name:  "sort-rows",
				Usage: helpers.DefaultSortFlag.Usage,
				Value: "value",
			},
			&cli.StringFlag{
				Name:  "sort-cols",
				Usage: helpers.DefaultSortFlag.Usage,
				Value: "value",
			},
			helpers.SnapshotFlag,
			helpers.NoOutFlag,
			helpers.CSVFlag,
			helpers.FormatFlag,
		},
	})
}
