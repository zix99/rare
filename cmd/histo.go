package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"rare/pkg/aggregation"
	"rare/pkg/multiterm"
	"sync"
	"time"

	"github.com/urfave/cli"
)

func writeOutput(writer *multiterm.HistoWriter, counter *aggregation.MatchCounter, count int) {
	items := counter.ItemsTop(count)
	for idx, match := range items {
		writer.WriteForLine(idx, match.Name, match.Item.Count())
	}
}

func histoFunction(c *cli.Context) error {
	topItems := c.Int("n")

	counter := aggregation.NewCounter()
	writer := multiterm.NewHistogram(topItems)
	writer.ShowBar = c.Bool("bars")
	done := make(chan bool)
	var mux sync.Mutex

	go func() {
		for {
			select {
			case <-done:
				return
			case <-time.After(50 * time.Millisecond):
				mux.Lock()
				writeOutput(writer, counter, topItems)
				mux.Unlock()
			}
		}
	}()

	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, os.Interrupt)
	extractor := buildExtractorFromArguments(c)
PROCESSING_LOOP:
	for {
		select {
		case <-exitSignal:
			break PROCESSING_LOOP
		case match, more := <-extractor.ReadChan:
			if !more {
				break PROCESSING_LOOP
			}
			mux.Lock()
			counter.Inc(match.Extracted)
			mux.Unlock()
		}
	}
	done <- true

	writeOutput(writer, counter, topItems)
	fmt.Println()

	writeExtractorSummary(extractor)
	fmt.Fprintf(os.Stderr, "Groups:  %d\n", counter.GroupCount())
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
		ArgsUsage: "<-|filename>",
		Flags: buildExtractorFlags(
			cli.BoolFlag{
				Name:  "bars,b",
				Usage: "Display bars as part of histogram",
			},
			cli.IntFlag{
				Name:  "num,n",
				Usage: "Number of elements to display",
				Value: 5,
			}),
	}
}
