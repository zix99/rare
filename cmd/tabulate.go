package cmd

import (
	"fmt"
	"os"
	. "rare/cmd/helpers"
	"rare/pkg/aggregation"
	"rare/pkg/color"
	"rare/pkg/humanize"
	"rare/pkg/multiterm"

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
	)

	counter := aggregation.NewTable(delim)
	writer := multiterm.NewTable(numCols, numRows)

	ext := BuildExtractorFromArguments(c)

	RunAggregationLoop(ext, counter, func() {
		cols := minColSlice(numCols, append([]string{""}, counter.OrderedColumns()...))
		writer.WriteRow(0, cols...)

		var rows []*aggregation.TableRow
		if sortByKeys {
			rows = counter.OrderedRowsByName()
		} else {
			rows = counter.OrderedRows()
		}

		for idx, row := range rows {
			rowVals := make([]string, len(cols)+1)
			rowVals[0] = row.Name()
			for idx, colName := range cols[1:] {
				rowVals[1+idx] = humanize.Hi(row.Value(colName))
			}
			writer.WriteRow(1+idx, rowVals...)
		}
	})

	fmt.Fprintf(os.Stderr, "Rows: %s; Cols: %s\n",
		color.Wrapi(color.Yellow, counter.RowCount()),
		color.Wrapi(color.BrightBlue, counter.ColumnCount()))

	return nil
}

func TabulateCommand() *cli.Command {
	return &cli.Command{
		Name:      "tabulate",
		Usage:     "Create a 2D summarizing table of extracted data",
		Action:    tabulateFunction,
		ArgsUsage: DefaultArgumentDescriptor,
		Flags: BuildExtractorFlags(
			cli.StringFlag{
				Name:  "delim",
				Usage: "Character to tabulate on. Use {tab} helper by default",
				Value: "\t",
			},
			cli.IntFlag{
				Name:  "num,n",
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
		),
	}
}
