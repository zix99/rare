package cmd

import (
	. "rare/cmd/helpers"
	"rare/pkg/aggregation"
	"rare/pkg/color"
	"rare/pkg/extractor"
	"rare/pkg/humanize"
	"rare/pkg/multiterm"
	"strconv"

	"github.com/urfave/cli"
)

func humanf(arg interface{}) string {
	return color.Wrap(color.BrightWhite, humanize.Hf(arg))
}

func writeAggrOutput(writer *multiterm.TermWriter, aggr *aggregation.MatchNumerical) {
	writer.WriteForLine(0, "Samples:  %v", color.Wrap(color.BrightWhite, humanize.Hi(aggr.Count())))
	writer.WriteForLine(1, "Mean:     %v", humanf(aggr.Mean()))
	writer.WriteForLine(2, "Min:      %v", humanf(aggr.Min()))
	writer.WriteForLine(3, "Max:      %v", humanf(aggr.Max()))

	writer.WriteForLine(4, "")

	data := aggr.Analyze()
	writer.WriteForLine(5, "Median:   %v", humanf(data.Median()))
	writer.WriteForLine(6, "Mode:     %v", humanf(data.Mode()))
	writer.WriteForLine(7, "P90:      %v", humanf(data.Quantile(0.9)))
	writer.WriteForLine(8, "P99:      %v", humanf(data.Quantile(0.99)))
	writer.WriteForLine(9, "P99.9:    %v", humanf(data.Quantile(0.999)))
}

func analyzeFunction(c *cli.Context) error {
	aggr := aggregation.NewNumericalAggregator()
	writer := multiterm.New(10)
	defer multiterm.ResetCursor()

	ext := BuildExtractorFromArguments(c)

	RunAggregationLoop(ext, func(match *extractor.Match) {
		val, err := strconv.ParseFloat(match.Extracted, 64)
		if err == nil {
			aggr.Sample(val)
		}
	}, func() {
		writeAggrOutput(writer, aggr)
	})

	return nil
}

func AnalyzeCommand() *cli.Command {
	return &cli.Command{
		Name:      "analyze",
		Usage:     "Numerical analysis on a set of filtered data",
		Action:    analyzeFunction,
		ArgsUsage: DefaultArgumentDescriptor,
		Flags: BuildExtractorFlags(
			cli.BoolFlag{
				Name:  "extra",
				Usage: "Displays extra analysis on the data (Requires more memory and cpu)",
			},
		),
	}
}
