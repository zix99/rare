package cmd

import (
	"fmt"

	"reblurb/pkg/aggregation"

	"github.com/urfave/cli"
)

func histoFunction(c *cli.Context) error {
	fmt.Println("Howdy")

	counter := aggregation.NewCounter()

	extractor := buildExtractorFromArguments(c)
	for {
		match, more := <-extractor.ReadChan
		if !more {
			break
		}
		counter.Inc(match.Extracted)
	}
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
