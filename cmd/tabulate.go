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
		delim   = c.String("delim")
		numRows = c.Int("num")
		numCols = c.Int("cols")
	)

	counter := aggregation.NewTable(delim)
	writer := multiterm.NewTable(numCols, numRows)

	ext := BuildExtractorFromArguments(c)

	RunAggregationLoop(ext, counter, func() {
		cols := minColSlice(numCols, counter.OrderedColumns())
		writer.WriteRow(0, cols...)
		for idx, row := range counter.OrderedRows() {
			rowVals := make([]string, len(cols)+1)
			rowVals[0] = row.Name()
			for idx, colName := range cols {
				rowVals[1+idx] = humanize.Hi(row.Value(colName))
			}
			writer.WriteRow(1+idx, rowVals...)
		}
	})

	if counter.ParseErrors() > 0 {
		fmt.Fprint(os.Stderr, color.Wrapf(color.Red, "Parse Errors: %v\n", humanize.Hi(counter.ParseErrors())))
	}

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
		),
	}
}
