package cmd

import (
	"fmt"
	"rare/cmd/helpers"
	"rare/pkg/aggregation"
	"rare/pkg/color"
	"rare/pkg/expressions"
	"rare/pkg/humanize"
	"rare/pkg/multiterm"
	"rare/pkg/multiterm/termrenderers"

	"github.com/urfave/cli"
)

func minColSlice(count int, cols []string) []string {
	if len(cols) < count {
		return cols
	}
	return cols[:count]
}

func tabulateFunction(c *cli.Context) error {
	var (
		delim      = c.String("delim")
		numRows    = c.Int("num")
		numCols    = c.Int("cols")
		sortByKeys = c.Bool("sortkey")
		rowtotals  = c.Bool("rowtotal") || c.Bool("x")
		coltotals  = c.Bool("coltotal") || c.Bool("x")
	)

	counter := aggregation.NewTable(delim)
	writer := termrenderers.NewTable(multiterm.New(), numCols+2, numRows+2)

	batcher := helpers.BuildBatcherFromArguments(c)
	ext := helpers.BuildExtractorFromArguments(c, batcher)

	helpers.RunAggregationLoop(ext, counter, func() {
		var cols []string
		if sortByKeys {
			cols = counter.OrderedColumnsByName()
		} else {
			cols = counter.OrderedColumns()
		}
		cols = minColSlice(numCols, cols) // Cap columns

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
		var rows []*aggregation.TableRow
		if sortByKeys {
			rows = counter.OrderedRowsByName()
		} else {
			rows = counter.OrderedRows()
		}

		line := 1
		for i := 0; i < len(rows) && i < numRows; i++ {
			row := rows[i]
			rowVals := make([]string, len(cols)+2)
			rowVals[0] = color.Wrap(color.Yellow, row.Name())
			for idx, colName := range cols {
				rowVals[idx+1] = humanize.Hi(row.Value(colName))
			}
			if rowtotals {
				rowVals[len(rowVals)-1] = color.Wrap(color.BrightBlack, humanize.Hi(row.Sum()))
			}
			writer.WriteRow(line, rowVals...)
			line++
		}

		// Write totals
		if coltotals {
			rowVals := make([]string, len(cols)+2)
			rowVals[0] = color.Wrap(color.BrightBlack+color.Underline, "Total")
			for idx, colName := range cols {
				rowVals[idx+1] = color.Wrap(color.BrightBlack, humanize.Hi(counter.ColTotal(colName)))
			}

			if rowtotals { // super total
				sum := counter.Sum()
				rowVals[len(rowVals)-1] = color.Wrap(color.BrightWhite, humanize.Hi(sum))
			}

			writer.WriteRow(line, rowVals...)
		}

		writer.WriteFooter(0, helpers.FWriteExtractorSummary(ext, counter.ParseErrors(),
			fmt.Sprintf("(R: %v; C: %v)", color.Wrapi(color.Yellow, counter.RowCount()), color.Wrapi(color.BrightBlue, counter.ColumnCount()))))
		writer.WriteFooter(1, batcher.StatusString())
	})

	writer.Close()

	return helpers.DetermineErrorState(batcher, ext, counter)
}

func tabulateCommand() *cli.Command {
	return helpers.AdaptCommandForExtractor(cli.Command{
		Name:      "tabulate",
		Aliases:   []string{"table"},
		ShortName: "t",
		Usage:     "Create a 2D summarizing table of extracted data",
		Description: `Summarizes the extracted data as a 2D data table.
		The expression key data format is {$ a b [c]}, where a is the column key,
		b is the rowkey, and optionally c is the increment value (Default: 1)`,
		Action: tabulateFunction,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "delim",
				Usage: "Character to tabulate on. Use {$} helper by default",
				Value: expressions.ArraySeparatorString,
			},
			cli.IntFlag{
				Name:  "num,n,rows",
				Usage: "Number of elements to display",
				Value: 20,
			},
			cli.IntFlag{
				Name:  "cols",
				Usage: "Number of columns to display",
				Value: 10,
			},
			cli.BoolFlag{
				Name:  "sortkey,sk",
				Usage: "Sort rows by key name rather than by values",
			},
			cli.BoolFlag{
				Name:  "rowtotal",
				Usage: "Show row totals",
			},
			cli.BoolFlag{
				Name:  "coltotal",
				Usage: "Show column totals",
			},
			cli.BoolFlag{
				Name:  "extra,x",
				Usage: "Display row and column totals",
			},
		},
	})
}
