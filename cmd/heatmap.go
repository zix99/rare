package cmd

import (
	"fmt"
	"rare/cmd/helpers"
	"rare/pkg/aggregation"
	"rare/pkg/color"
	"rare/pkg/expressions"
	"rare/pkg/multiterm"
	"rare/pkg/multiterm/termrenderers"

	"github.com/urfave/cli/v2"
)

func heatmapFunction(c *cli.Context) error {
	var (
		delim    = c.String("delim")
		numRows  = c.Int("num")
		numCols  = c.Int("cols")
		minFixed = c.IsSet("min")
		minVal   = c.Int64("min")
		maxFixed = c.IsSet("max")
		maxVal   = c.Int64("max")
		sortRows = c.String("sort-rows")
		sortCols = c.String("sort-cols")
	)

	counter := aggregation.NewTable(delim)

	batcher := helpers.BuildBatcherFromArguments(c)
	ext := helpers.BuildExtractorFromArguments(c, batcher)
	rowSorter := helpers.BuildSorterOrFail(sortRows)
	colSorter := helpers.BuildSorterOrFail(sortCols)

	vt := helpers.BuildVTermFromArguments(c)
	writer := termrenderers.NewHeatmap(vt, numRows, numCols)

	writer.FixedMin = minFixed
	writer.FixedMax = maxFixed
	if minFixed || maxFixed {
		writer.UpdateMinMax(minVal, maxVal)
	}

	helpers.RunAggregationLoop(ext, counter, func() {
		writer.WriteTable(counter, rowSorter, colSorter)
		writer.WriteFooter(0, helpers.FWriteExtractorSummary(ext, counter.ParseErrors(),
			fmt.Sprintf("(R: %v; C: %v)", color.Wrapi(color.Yellow, counter.RowCount()), color.Wrapi(color.BrightBlue, counter.ColumnCount()))))
		writer.WriteFooter(1, batcher.StatusString())
	})

	writer.Close()

	return helpers.DetermineErrorState(batcher, ext, counter)
}

func heatmapCommand() *cli.Command {
	return helpers.AdaptCommandForExtractor(cli.Command{
		Name:    "heatmap",
		Aliases: []string{"heat", "hm"},
		Usage:   "Create a 2D heatmap of extracted data",
		Description: `Creates a dense 2D visual of extracted data.  Each character
		represents a single data-point, and can create an alternative visualization to
		a table.  Unicode and color support required for effective display`,
		Action: heatmapFunction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "delim",
				Usage: "Character to tabulate on. Use {$} helper by default",
				Value: expressions.ArraySeparatorString,
			},
			&cli.IntFlag{
				Name:    "num",
				Aliases: []string{"rows", "n"},
				Usage:   "Number of elements (rows) to display",
				Value:   20,
			},
			&cli.IntFlag{
				Name:  "cols",
				Usage: "Number of columns to display",
				Value: multiterm.TermCols() - 15,
			},
			&cli.Int64Flag{
				Name:  "min",
				Usage: "Sets the lower bounds of the heatmap (default: auto)",
			},
			&cli.Int64Flag{
				Name:  "max",
				Usage: "Sets the upper bounds of the heatmap (default: auto)",
			},
			&cli.StringFlag{
				Name:  "sort-rows",
				Usage: helpers.DefaultSortFlag.Usage,
				Value: helpers.DefaultSortFlag.Value,
			},
			&cli.StringFlag{
				Name:  "sort-cols",
				Usage: helpers.DefaultSortFlag.Usage,
				Value: helpers.DefaultSortFlag.Value,
			},
			helpers.SnapshotFlag,
		},
	})
}
