package cmd

import (
	"fmt"
	"rare/cmd/helpers"
	"rare/cmd/readProgress"
	"rare/pkg/aggregation"
	"rare/pkg/color"
	"rare/pkg/multiterm"

	"github.com/urfave/cli"
)

func distFunction(c *cli.Context) error {
	var (
		topItems    = c.Int("n")
		reverseSort = c.Bool("reverse")
		atLeast     = c.Int64("atleast")
		sortByKey   = c.Bool("sk")
		extra       = c.Bool("extra")
	)

	counter := aggregation.NewFuzzyAggregator(0.8)
	writer := multiterm.NewHistogram(multiterm.New(), 10)
	writer.ShowBar = c.Bool("bars") || extra
	writer.ShowPercentage = c.Bool("percentage") || extra

	ext := helpers.BuildExtractorFromArguments(c)

	helpers.RunAggregationLoop(ext, counter, func() {
		writeHistoOutput(writer, counter.Histo, topItems, reverseSort, sortByKey, atLeast)
		writer.InnerWriter().WriteForLine(topItems, helpers.FWriteExtractorSummary(ext,
			counter.ParseErrors(),
			fmt.Sprintf("(Groups: %s)", color.Wrapi(color.BrightBlue, counter.Histo.GroupCount()))))
		writer.InnerWriter().WriteForLine(topItems+1, readProgress.GetReadFileString())
	})

	writer.InnerWriter().Close()

	return nil
}

func distCommand() *cli.Command {
	return helpers.AdaptCommandForExtractor(cli.Command{
		Name:      "distance",
		ShortName: "d",
		Aliases:   []string{"dist"},
		Description: `Generates a live-updating histogram of the input data, looking
		for a relative distance between various results.  This is useful to find
		similar log messages that may have slight differences to them (eg ids)
		and aggregating and search for these messages`,
		Action: distFunction,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "bars,b",
				Usage: "Display bars as part of histogram",
			},
			cli.BoolFlag{
				Name:  "percentage",
				Usage: "Display percentage of total next to the value",
			},
			cli.BoolFlag{
				Name:  "extra,x",
				Usage: "Alias for -b --percentage",
			},
			cli.IntFlag{
				Name:  "num,n",
				Usage: "Number of elements to display",
				Value: 5,
			},
			cli.Int64Flag{
				Name:  "atleast",
				Usage: "Only show results if there are at least this many samples",
				Value: 0,
			},
			cli.BoolFlag{
				Name:  "reverse",
				Usage: "Reverses the display sort-order",
			},
			cli.BoolFlag{
				Name:  "sortkey,sk",
				Usage: "Sort by key, rather than value",
			},
		},
	})
}
