package cmd

import (
	"rare/cmd/helpers"
	"rare/pkg/aggregation"
	"rare/pkg/expressions"
	"rare/pkg/multiterm"
	"rare/pkg/multiterm/termrenderers"

	"github.com/urfave/cli"
)

func heatmapFunction(c *cli.Context) error {
	counter := aggregation.NewTable(expressions.ArraySeparatorString)

	batcher := helpers.BuildBatcherFromArguments(c)
	ext := helpers.BuildExtractorFromArguments(c, batcher)

	writer := termrenderers.NewHeatmap(multiterm.New(), 20, 10) // TODO: Configurable size

	helpers.RunAggregationLoop(ext, counter, func() {
		writer.WriteTable(counter)
		writer.WriteFooter(0, helpers.FWriteExtractorSummary(ext, counter.ParseErrors()))
		writer.WriteFooter(1, batcher.StatusString())
	})

	return helpers.DetermineErrorState(batcher, ext, nil)
}

func heatmapCommand() *cli.Command {
	return helpers.AdaptCommandForExtractor(cli.Command{
		Name:   "heatmap",
		Action: heatmapFunction,
	})
}
