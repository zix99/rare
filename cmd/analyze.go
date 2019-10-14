package cmd

import (
	"log"
	. "rare/cmd/helpers"
	"rare/pkg/aggregation"
	"rare/pkg/color"
	"rare/pkg/humanize"
	"rare/pkg/multiterm"
	"strconv"

	"github.com/urfave/cli"
)

func humanf(arg interface{}) string {
	return color.Wrap(color.BrightWhite, humanize.Hf(arg))
}

func writeAggrOutput(writer *multiterm.TermWriter, aggr *aggregation.MatchNumerical, extra bool, quantiles []float64) {
	writer.WriteForLine(0, "Samples:  %v", color.Wrap(color.BrightWhite, humanize.Hi(aggr.Count())))
	writer.WriteForLine(1, "Mean:     %v", humanf(aggr.Mean()))
	writer.WriteForLine(2, "Min:      %v", humanf(aggr.Min()))
	writer.WriteForLine(3, "Max:      %v", humanf(aggr.Max()))

	if extra {
		writer.WriteForLine(4, "")

		data := aggr.Analyze()
		writer.WriteForLine(5, "Median:   %v", humanf(data.Median()))
		writer.WriteForLine(6, "Mode:     %v", humanf(data.Mode()))
		writer.WriteForLine(7, "StdDev:   %v", humanf(aggr.StdDev()))
		for idx, q := range quantiles {
			writer.WriteForLine(8+idx, "P%02.4f: %v", q, humanf(data.Quantile(q/100.0)))
		}
	}
}

func parseStringSet(vals []string) []float64 {
	ret := make([]float64, len(vals))
	for idx, val := range vals {
		parsedVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			log.Fatalf("%s is not a number: %v", val, err)
		}
		ret[idx] = parsedVal
	}
	return ret
}

func analyzeFunction(c *cli.Context) error {
	config := aggregation.NumericalConfig{
		Reverse: c.Bool("reverse"),
	}
	extra := c.Bool("extra")
	quantiles := parseStringSet(c.StringSlice("quantile"))

	aggr := aggregation.NewNumericalAggregator(&config)
	writer := multiterm.New(15)
	defer multiterm.ResetCursor()

	ext := BuildExtractorFromArguments(c)

	RunAggregationLoop(ext, aggr, func() {
		writeAggrOutput(writer, aggr, extra, quantiles)
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
			cli.BoolFlag{
				Name:  "reverse,r",
				Usage: "Reverses the numerical series when ordered-analysis takes place (eg Quantile)",
			},
			cli.StringSliceFlag{
				Name:  "quantile,q",
				Usage: "Adds a quantile to the output set. Requires --extra",
				Value: &cli.StringSlice{"90", "99", "99.9"},
			},
		),
	}
}
