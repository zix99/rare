package cmd

import (
	"fmt"
	"os"
	. "rare/cmd/helpers"
	"rare/pkg/aggregation"
	"rare/pkg/color"
	"rare/pkg/humanize"
	"rare/pkg/multiterm"
	"strings"

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
	writer := multiterm.New(numRows)

	ext := BuildExtractorFromArguments(c)

	RunAggregationLoop(ext, counter, func() {
		cols := minColSlice(numCols, counter.Columns())
		writer.WriteForLine(0, "name\t"+strings.Join(cols, "\t"))
		for idx, row := range counter.OrderedRows() {
			var sb strings.Builder
			sb.WriteString(row.Name())
			for _, colName := range cols {
				sb.WriteRune('\t')
				sb.WriteString(humanize.Hi(row.Value(colName)))
			}
			writer.WriteForLine(1+idx, sb.String())
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
