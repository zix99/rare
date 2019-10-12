package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"rare/pkg/aggregation"
	"rare/pkg/color"
	"rare/pkg/multiterm"
	"sync"
	"sync/atomic"
	"time"

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
	done := make(chan bool)

	var mux sync.Mutex
	var hasUpdates atomic.Value
	hasUpdates.Store(false)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-time.After(100 * time.Millisecond):
				if hasUpdates.Load().(bool) {
					hasUpdates.Store(false)
					mux.Lock()
					writeOutput(writer, counter, topItems, reverseSort, sortByKey)
					mux.Unlock()
				}
			}
		}
	}()

	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, os.Interrupt)
	extractor := buildExtractorFromArguments(c)
	readChan := extractor.ReadChan()
PROCESSING_LOOP:
	for {
		select {
		case <-exitSignal:
			break PROCESSING_LOOP
		case match, more := <-readChan:
			if !more {
				break PROCESSING_LOOP
			}
			mux.Lock()
			counter.Inc(match.Extracted)
			hasUpdates.Store(true)
			mux.Unlock()
		}
	}
	done <- true

	writeOutput(writer, counter, topItems, reverseSort, sortByKey)
	fmt.Println()

	writeExtractorSummary(extractor)
	fmt.Fprintf(os.Stderr, "Groups:  %s\n", color.Wrapf(color.BrightWhite, "%d", counter.GroupCount()))
	multiterm.ResetCursor()

	return nil
}

// HistogramCommand Exported command
func HistogramCommand() *cli.Command {
	return &cli.Command{
		Name:      "histo",
		Usage:     "Summarize results by extracting them to a histogram",
		Action:    histoFunction,
		Aliases:   []string{"histogram"},
		ShortName: "h",
		ArgsUsage: "<-|filename|glob...>",
		Flags: buildExtractorFlags(
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
			}),
	}
}
