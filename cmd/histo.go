package cmd

import (
	"fmt"
	. "rare/cmd/helpers"
	"rare/pkg/aggregation"
	"rare/pkg/color"
	"rare/pkg/multiterm"

	"github.com/urfave/cli"
)

func writeHistoOutput(writer *multiterm.HistoWriter, counter *aggregation.MatchCounter, count int, reverse bool, sortByKey bool) {
	var items []aggregation.MatchPair
	if sortByKey {
		items = counter.ItemsSortedByKey(count, reverse)
	} else {
		items = counter.ItemsSorted(count, reverse)
	}
	for idx, match := range items {
		writer.WriteForLine(idx, match.Name, match.Item.Count())
	}
}

func histoFunction(c *cli.Context) error {
	var (
		topItems    = c.Int("n")
		reverseSort = c.Bool("reverse")
		sortByKey   = c.Bool("sk")
	)

	counter := aggregation.NewCounter()
	writer := multiterm.NewHistogram(multiterm.New(), topItems)
	writer.ShowBar = c.Bool("bars")

	ext := BuildExtractorFromArguments(c)

	RunAggregationLoop(ext, counter, func() {
		writeHistoOutput(writer, counter, topItems, reverseSort, sortByKey)
		writer.InnerWriter().WriteForLine(topItems, FWriteExtractorSummary(ext,
			counter.ParseErrors(),
			fmt.Sprintf(" (Groups: %s)", color.Wrapi(color.BrightBlue, counter.GroupCount()))))
		writer.InnerWriter().WriteForLine(topItems+1, GetReadFileString())
	})

	writer.InnerWriter().Close()

	return nil
}

// HistogramCommand Exported command
func histogramCommand() *cli.Command {
	return AdaptCommandForExtractor(cli.Command{
		Name:  "histogram",
		Usage: "Summarize results by extracting them to a histogram",
		Description: `Generates a live-updating histogram of the extracted information from a file
		Each line in the file will be matched, any the matching part extracted
		as a key and counted.
		If an extraction expression is provided with -e, that will be used
		as the key instead`,
		Action:    histoFunction,
		Aliases:   []string{"histo"},
		ShortName: "h",
		ArgsUsage: DefaultArgumentDescriptor,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "bars,b",
				Usage: "Display bars as part of histogram",
			},
			cli.IntFlag{
				Name:  "num,n",
				Usage: "Number of elements to display",
				Value: 5,
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
