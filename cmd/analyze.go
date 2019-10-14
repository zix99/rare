package cmd

import (
	"fmt"
	. "rare/cmd/helpers"
	"rare/pkg/aggregation"
	"rare/pkg/multiterm"
	"strconv"
	"sync"
	"time"

	"github.com/urfave/cli"
)

func writeAggrOutput(writer *multiterm.TermWriter, aggr *aggregation.MatchNumerical) {
	writer.WriteForLine(0, "Samples: %d", aggr.Count())
	writer.WriteForLine(1, "Mean:    %f", aggr.Mean())
	writer.WriteForLine(2, "Min:     %f", aggr.Min())
	writer.WriteForLine(3, "Max:     %f", aggr.Max())

	data := aggr.Analyze()
	writer.WriteForLine(5, "Median:   %f", data.Median())
	writer.WriteForLine(6, "Mode:     %f", data.Mode())
	writer.WriteForLine(7, "P90:      %f", data.Quantile(0.9))
	writer.WriteForLine(8, "P99:      %f", data.Quantile(0.99))
	writer.WriteForLine(9, "P99.9:    %f", data.Quantile(0.999))
}

func analyzeFunction(c *cli.Context) error {
	aggr := aggregation.NewNumericalAggregator()
	writer := multiterm.New(10)
	defer multiterm.ResetCursor()

	var mux sync.Mutex

	extractor := BuildExtractorFromArguments(c)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-time.After(100 * time.Millisecond):
				mux.Lock()
				writeAggrOutput(writer, aggr)
				mux.Unlock()
			}
		}
	}()

	readChan := extractor.ReadChan()
	for {
		match, more := <-readChan
		if !more {
			break
		}

		val, err := strconv.ParseFloat(match.Extracted, 64)
		if err == nil {
			mux.Lock()
			aggr.Sample(val)
			mux.Unlock()
		}
	}
	done <- true

	writeAggrOutput(writer, aggr)
	fmt.Println()

	WriteExtractorSummary(extractor)

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
