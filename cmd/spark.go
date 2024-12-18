package cmd

import (
	"fmt"
	"rare/cmd/helpers"
	"rare/pkg/aggregation"
	"rare/pkg/color"
	"rare/pkg/csv"
	"rare/pkg/expressions"
	"rare/pkg/multiterm"
	"rare/pkg/multiterm/termrenderers"

	"github.com/urfave/cli/v2"
)

func sparkFunction(c *cli.Context) error {
	var (
		delim      = c.String("delim")
		numRows    = c.Int("num")
		numCols    = c.Int("cols")
		noTruncate = c.Bool("notruncate")
		scalerName = c.String(helpers.ScaleFlag.Name)
		sortRows   = c.String("sort-rows")
		sortCols   = c.String("sort-cols")
	)

	counter := aggregation.NewTable(delim)

	batcher := helpers.BuildBatcherFromArguments(c)
	ext := helpers.BuildExtractorFromArguments(c, batcher)
	rowSorter := helpers.BuildSorterOrFail(sortRows)
	colSorter := helpers.BuildSorterOrFail(sortCols)

	vt := helpers.BuildVTermFromArguments(c)
	writer := termrenderers.NewSpark(vt, numRows, numCols)
	writer.Scaler = helpers.BuildScalerOrFail(scalerName)

	helpers.RunAggregationLoop(ext, counter, func() {

		// Trim unused data from the data store (keep memory tidy!)
		if !noTruncate {
			if keepCols := counter.OrderedColumns(colSorter); len(keepCols) > numCols {
				keepCols = keepCols[len(keepCols)-numCols:]
				keepLookup := make(map[string]struct{})
				for _, item := range keepCols {
					keepLookup[item] = struct{}{}
				}
				counter.Trim(func(col, row string, val int64) bool {
					_, ok := keepLookup[col]
					return !ok
				})
			}
		}

		// Write spark
		writer.WriteTable(counter, rowSorter, colSorter)
		writer.WriteFooter(0, helpers.FWriteExtractorSummary(ext, counter.ParseErrors(),
			fmt.Sprintf("(R: %v; C: %v)", color.Wrapi(color.Yellow, counter.RowCount()), color.Wrapi(color.BrightBlue, counter.ColumnCount()))))
		writer.WriteFooter(1, batcher.StatusString())
	})

	// Not deferred intentionally
	writer.Close()

	if err := helpers.TryWriteCSV(c, counter, csv.WriteTable); err != nil {
		return err
	}

	return helpers.DetermineErrorState(batcher, ext, counter)
}

func sparkCommand() *cli.Command {
	return helpers.AdaptCommandForExtractor(cli.Command{
		Name:    "spark",
		Aliases: []string{"sparkline", "s"},
		Usage:   "Create rows of sparkline graphs",
		Description: `Create rows of a sparkkline graph, all scaled equally
		based on a table like input`,
		Category: cmdCatVisualize,
		Action:   sparkFunction,
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
			&cli.BoolFlag{
				Name:  "notruncate",
				Usage: "Disable truncating data that doesnt fit in the sparkline",
				Value: false,
			},
			&cli.StringFlag{
				Name:  "sort-rows",
				Usage: helpers.DefaultSortFlag.Usage,
				Value: "value",
			},
			&cli.StringFlag{
				Name:  "sort-cols",
				Usage: helpers.DefaultSortFlag.Usage,
				Value: "numeric",
			},
			helpers.SnapshotFlag,
			helpers.NoOutFlag,
			helpers.CSVFlag,
			helpers.ScaleFlag,
		},
	})
}
