package cmd

import (
	"fmt"
	"os"
	"rare/cmd/helpers"
	"rare/pkg/aggregation"
	"rare/pkg/aggregation/sorting"
	"rare/pkg/color"
	"rare/pkg/multiterm"
	"rare/pkg/multiterm/termrenderers"

	"github.com/urfave/cli/v2"
)

func writeHistoOutput(writer *termrenderers.HistoWriter, counter *aggregation.MatchCounter, count int, sorter sorting.NameValueSorter, atLeast int64) {
	items := counter.ItemsSortedBy(count, sorter)
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
		reverseSort = c.Bool("sort-reverse")
		sortName    = c.String("sort")
		atLeast     = c.Int64("atleast")
		extra       = c.Bool("extra")
		all         = c.Bool("all")
	)

	counter := aggregation.NewCounter()
	writer := termrenderers.NewHistogram(multiterm.New(), topItems)
	writer.ShowBar = c.Bool("bars") || extra
	writer.ShowPercentage = c.Bool("percentage") || extra

	batcher := helpers.BuildBatcherFromArguments(c)
	ext := helpers.BuildExtractorFromArguments(c, batcher)
	sorter := helpers.BuildSorter(sortName, reverseSort)

	progressString := func() string {
		return helpers.FWriteExtractorSummary(ext,
			counter.ParseErrors(),
			fmt.Sprintf("(Groups: %s)", color.Wrapi(color.BrightBlue, counter.GroupCount())))
	}

	helpers.RunAggregationLoop(ext, counter, func() {
		writeHistoOutput(writer, counter, topItems, sorter, atLeast)
		writer.WriteFooter(0, progressString())
		writer.WriteFooter(1, batcher.StatusString())
	})

	writer.Close()

	if all {
		fmt.Println("Full Table:")
		vterm := multiterm.NewVirtualTerm()
		vWriter := termrenderers.NewHistogram(vterm, counter.GroupCount())
		writeHistoOutput(vWriter, counter, counter.GroupCount(), sorter, atLeast)

		vterm.WriteToOutput(os.Stdout)
		fmt.Println(progressString())
	}

	return helpers.DetermineErrorState(batcher, ext, counter)
}

// HistogramCommand Exported command
func histogramCommand() *cli.Command {
	return helpers.AdaptCommandForExtractor(cli.Command{
		Name:  "histogram",
		Usage: "Summarize results by extracting them to a histogram",
		Description: `Generates a live-updating histogram of the extracted information from a file
		Each line in the file will be matched, any the matching part extracted
		as a key and counted.
		If an extraction expression is provided with -e, that will be used
		as the key instead
		If multiple values are provided via the array syntax {$} or multiple expressions,
		then the 2nd value will be used as the count incrementor`,
		Action:    histoFunction,
		Aliases:   []string{"histo", "h"},
		ArgsUsage: helpers.DefaultArgumentDescriptor,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "After summarization is complete, print all histogram buckets",
			},
			&cli.BoolFlag{
				Name:    "bars",
				Aliases: []string{"b"},
				Usage:   "Display bars as part of histogram",
			},
			&cli.BoolFlag{
				Name:  "percentage",
				Usage: "Display percentage of total next to the value",
			},
			&cli.BoolFlag{
				Name:    "extra",
				Aliases: []string{"x"},
				Usage:   "Alias for -b --percentage",
			},
			&cli.IntFlag{
				Name:    "num",
				Aliases: []string{"n"},
				Usage:   "Number of elements to display",
				Value:   5,
			},
			&cli.Int64Flag{
				Name:  "atleast",
				Usage: "Only show results if there are at least this many samples",
				Value: 0,
			},
			&cli.StringFlag{
				Name:  "sort",
				Usage: "Sets sorting method (value, text, numeric, contextual)",
				Value: "value",
			},
			&cli.BoolFlag{
				Name:    "sort-reverse",
				Aliases: []string{"reverse"},
				Usage:   "Reverses the display sort-order",
			},
		},
	})
}
