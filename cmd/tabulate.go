package cmd

import (
	"fmt"

	"github.com/zix99/rare/cmd/helpers"
	"github.com/zix99/rare/pkg/aggregation"
	"github.com/zix99/rare/pkg/color"
	"github.com/zix99/rare/pkg/csv"
	"github.com/zix99/rare/pkg/expressions"
	"github.com/zix99/rare/pkg/multiterm/termrenderers"

	"github.com/urfave/cli/v2"
)

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
	writer := termrenderers.NewDataTable(vt, numCols, numRows)
	writer.ShowRowTotals = rowtotals
	writer.ShowColTotals = coltotals

	batcher := helpers.BuildBatcherFromArguments(c)
	ext := helpers.BuildExtractorFromArguments(c, batcher)
	rowSorter := helpers.BuildSorterOrFail(sortRows)
	colSorter := helpers.BuildSorterOrFail(sortCols)

	if formatExp != "" {
		writer.SetFormatter(helpers.BuildFormatterOrFail(formatExp))
	}

	interrupt := helpers.RunAggregationLoop(ext, counter, func() {
		writer.WriteTable(counter, rowSorter, colSorter)

		writer.WriteFooter(0, helpers.BuildExtractorSummary(ext, counter.ParseErrors(),
			fmt.Sprintf("(R: %v; C: %v)", color.Wrapi(color.Yellow, counter.RowCount()), color.Wrapi(color.BrightBlue, counter.ColumnCount()))))
		writer.WriteFooter(1, batcher.StatusString())
	})

	writer.Close()

	if err := helpers.TryWriteCSV(c, counter, csv.WriteTable); err != nil {
		return err
	}

	return helpers.DetermineErrorState(interrupt, batcher, ext, counter)
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
