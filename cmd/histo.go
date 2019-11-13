package cmd

import (
	"fmt"
	"os"
	. "rare/cmd/helpers"
	"rare/pkg/aggregation"
	"rare/pkg/color"
	"rare/pkg/multiterm"

	"github.com/urfave/cli"
)

func writeOutput(writer *multiterm.HistoWriter, counter *aggregation.MatchCounter, count int, reverse bool, sortByKey bool) {
	var items []aggregation.MatchPair
	if sortByKey {
		items = counter.ItemsSortedByKey(count, reverse)
	} else {
		items = counter.ItemsSorted(count, reverse)
	}
	for idx, match := range items {
		writer.WriteForLine(idx, match.Name, match.Item.Count())
	}
	writer.InnerWriter().GoToBottom(0)
	writer.InnerWriter().WriteAtCursor(GetReadFileString())
}

func histoFunction(c *cli.Context) error {
	var (
		topItems    = c.Int("n")
		reverseSort = c.Bool("reverse")
		sortByKey   = c.Bool("sk")
	)

	counter := aggregation.NewCounter()
	writer := multiterm.NewHistogram(topItems)
	writer.ShowBar = c.Bool("bars")

	ext := BuildExtractorFromArguments(c)

	RunAggregationLoop(ext, counter, func() {
		writeOutput(writer, counter, topItems, reverseSort, sortByKey)
	})

	fmt.Fprintf(os.Stderr, "Groups:  %s\n", color.Wrapf(color.BrightWhite, "%d", counter.GroupCount()))

	return nil
}

// HistogramCommand Exported command
func histogramCommand() *cli.Command {
	return AdaptCommandForExtractor(cli.Command{
		Name:      "histogram",
		Usage:     "Summarize results by extracting them to a histogram",
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
