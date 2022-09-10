package cmd

import (
	"rare/cmd/helpers"
	"rare/pkg/aggregation"
	"rare/pkg/multiterm"
	"rare/pkg/multiterm/termrenderers"

	"github.com/urfave/cli/v2"
)

/*
Test Command:
go run . bars -sz -m "\[(.+?)\].*\" (\d+)" -e "{$ {buckettime {1} year nginx} {2}}" testdata/*
*/

func bargraphFunction(c *cli.Context) error {
	var (
		stacked  = c.Bool("stacked")
		sortName = c.String(helpers.DefaultSortFlag.Name)
	)

	counter := aggregation.NewSubKeyCounter()
	writer := termrenderers.NewBarGraph(multiterm.New())
	writer.Stacked = stacked

	batcher := helpers.BuildBatcherFromArguments(c)
	ext := helpers.BuildExtractorFromArguments(c, batcher)
	sorter := helpers.BuildSorterOrFail(sortName)

	helpers.RunAggregationLoop(ext, counter, func() {
		line := 0

		writer.SetKeys(counter.SubKeys()...)
		for _, row := range counter.ItemsSorted(sorter) {
			writer.WriteBar(line, row.Name, row.Item.Items()...)
			line++
		}

		writer.WriteFooter(0, helpers.FWriteExtractorSummary(ext, counter.ParseErrors()))
		writer.WriteFooter(1, batcher.StatusString())
	})

	writer.Close()

	return helpers.DetermineErrorState(batcher, ext, counter)
}

func bargraphCommand() *cli.Command {
	return helpers.AdaptCommandForExtractor(cli.Command{
		Name:    "bargraph",
		Aliases: []string{"bars", "bar", "b"},
		Usage:   "Create a bargraph of the given 1 or 2 dimension data",
		Description: `Creates a bargraph of one or two dimensional data.  Unlike histogram
		the bargraph can collapse and stack data in different formats.  The key data format
		is {$ a b [c]}, where a is the base-key, b is the optional sub-key, and c is the increment
		(defaults to 1)`,
		Action: bargraphFunction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "stacked",
				Aliases: []string{"s"},
				Usage:   "Display bargraph as stacked",
			},
			helpers.DefaultSortFlag,
		},
	})
}
