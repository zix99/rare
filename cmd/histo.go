package cmd

import (
	"fmt"
	"os"
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
	const topItems = 5
	counter := aggregation.NewCounter()
	writer := multiterm.NewHistogram(topItems)
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

	extractor := buildExtractorFromArguments(c)
	for {
		match, more := <-extractor.ReadChan
		if !more {
			break
		}
		mux.Lock()
		counter.Inc(match.Extracted)
		mux.Unlock()
	}

	writeOutput(writer, counter, topItems)
	fmt.Println()
	writeExtractorSummary(extractor)
	fmt.Fprintf(os.Stderr, "Groups:  %d\n", counter.GroupCount())
	multiterm.ResetCursor()

	done <- true
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
		Flags:     buildExtractorFlags(),
	}
}
