package cmd

import (
	"fmt"
	"reblurb/pkg/aggregation"
	"reblurb/pkg/multiterm"
	"time"

	"github.com/urfave/cli"
)

func writeOutput(writer *multiterm.TermWriter, counter *aggregation.MatchCounter, count int) {
	items := counter.ItemsTop(count)
	for idx, match := range items {
		writer.WriteForLine(idx, "%-32v %d", match.Name, match.Item.Count())
	}
}

func histoFunction(c *cli.Context) error {
	const topItems = 5
	counter := aggregation.NewCounter()
	writer := multiterm.New(topItems)
	done := make(chan bool)

	// TODO: Async safety
	go func() {
		for {
			select {
			case <-done:
				return
			case <-time.After(50 * time.Millisecond):
				writeOutput(writer, counter, topItems)
			}
		}
	}()

	extractor := buildExtractorFromArguments(c)
	for {
		match, more := <-extractor.ReadChan
		if !more {
			break
		}
		counter.Inc(match.Extracted)
	}
	writeOutput(writer, counter, topItems)
	fmt.Println()
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
		Flags:     []cli.Flag{},
	}
}
