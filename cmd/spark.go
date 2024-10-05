package cmd

import (
	"fmt"
	"rare/cmd/helpers"
	"rare/pkg/aggregation"
	"rare/pkg/aggregation/sorting"
	"rare/pkg/color"
	"rare/pkg/expressions"
	"rare/pkg/multiterm/termrenderers"
	"rare/pkg/multiterm/termscaler"

	"github.com/urfave/cli/v2"
)

func sparkFunction(c *cli.Context) error {
	// TODO: Truncate/trim table flag (or no-truncate flag?)
	// TODO: Table decoration flags (eg first/last data point, average?)

	counter := aggregation.NewTable(expressions.ArraySeparatorString) // TODO: argument

	batcher := helpers.BuildBatcherFromArguments(c)
	ext := helpers.BuildExtractorFromArguments(c, batcher)
	rowSorter := sorting.NVNameSorter // TODO
	colSorter := sorting.NVNameSorter // TODO

	vt := helpers.BuildVTermFromArguments(c)
	writer := termrenderers.NewSpark(vt, 10, 30) // TODO: Args

	writer.Scaler = termscaler.ScalerLinear // TODO

	helpers.RunAggregationLoop(ext, counter, func() {
		writer.WriteTable(counter, rowSorter, colSorter)
		writer.WriteFooter(0, helpers.FWriteExtractorSummary(ext, counter.ParseErrors(),
			fmt.Sprintf("(R: %v; C: %v)", color.Wrapi(color.Yellow, counter.RowCount()), color.Wrapi(color.BrightBlue, counter.ColumnCount()))))
		writer.WriteFooter(1, batcher.StatusString())
	})

	writer.Close()

	// TODO: Csv?

	return nil
}

func sparkCommand() *cli.Command {
	return helpers.AdaptCommandForExtractor(cli.Command{
		Name:   "spark",
		Action: sparkFunction,
	})
}
