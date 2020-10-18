package cmd

import (
	"fmt"
	"os"
	. "rare/cmd/helpers" //lint:ignore ST1001 Legacy
	"rare/cmd/readProgress"
	"rare/pkg/aggregation"
	"rare/pkg/color"
	"rare/pkg/multiterm"

	"github.com/urfave/cli"
)

func writeHistoOutput(writer *multiterm.HistoWriter, counter *aggregation.MatchCounter, count int, reverse bool, sortByKey bool, atLeast int64) {
	var items []aggregation.MatchPair
	if sortByKey {
		items = counter.ItemsSortedByKey(count, reverse)
	} else {
		items = counter.ItemsSorted(count, reverse)
	}
	line := 0
	writer.UpdateSamples(counter.Count())
	for _, match := range items {
		count := match.Item.Count()
		if count >= atLeast {
			writer.WriteForLine(line, match.Name, count)
			line++
		}
	}
}

func histoFunction(c *cli.Context) error {
	var (
		topItems    = c.Int("n")
		reverseSort = c.Bool("reverse")
		sortByKey   = c.Bool("sk")
		atLeast     = c.Int64("atleast")
		extra       = c.Bool("extra")
		all         = c.Bool("all")
	)

	counter := aggregation.NewCounter()
	writer := multiterm.NewHistogram(multiterm.New(), topItems)
	writer.ShowBar = c.Bool("bars") || extra
	writer.ShowPercentage = c.Bool("percentage") || extra

	ext := BuildExtractorFromArguments(c)

	progressString := func() string {
		return FWriteExtractorSummary(ext,
			counter.ParseErrors(),
			fmt.Sprintf("(Groups: %s)", color.Wrapi(color.BrightBlue, counter.GroupCount())))
	}

	RunAggregationLoop(ext, counter, func() {
		writeHistoOutput(writer, counter, topItems, reverseSort, sortByKey, atLeast)
		writer.InnerWriter().WriteForLine(topItems, progressString())
		writer.InnerWriter().WriteForLine(topItems+1, readProgress.GetReadFileString())
	})

	writer.InnerWriter().Close()

	if all {
		fmt.Println("Full Table:")
		vterm := multiterm.NewVirtualTerm()
		vWriter := multiterm.NewHistogram(vterm, counter.GroupCount())
		writeHistoOutput(vWriter, counter, counter.GroupCount(), reverseSort, sortByKey, atLeast)

		vterm.WriteToOutput(os.Stdout)
		fmt.Println(progressString())
	}

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
		as the key instead
		If multiple values are provided via the array syntax {$}, then the
		2nd value will be used as the count incrementor`,
		Action:    histoFunction,
		Aliases:   []string{"histo"},
		ShortName: "h",
		ArgsUsage: DefaultArgumentDescriptor,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "all,a",
				Usage: "After summarization is complete, print all histogram buckets",
			},
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
