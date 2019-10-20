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

func tabulateFunction(c *cli.Context) error {
	var (
		delim = c.String("delim")
	)

	counter := aggregation.NewTable(delim)
	writer := multiterm.New(10)

	ext := BuildExtractorFromArguments(c)

	RunAggregationLoop(ext, counter, func() {
		cols := counter.Columns()
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
				Usage: "Character to tabulate on",
				Value: " ",
			},
		),
	}
}
