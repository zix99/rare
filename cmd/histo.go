package cmd

import (
	"fmt"

	"github.com/urfave/cli"
)

func histoFunction(c *cli.Context) error {
	fmt.Println("Howdy")

	format := c.String("extract")

	extractor := buildExtractorFromArguments(c)
	for {
		match, more := <-extractor.ReadChan
		if !more {
			break
		}
		if format != "" {
			//fmt.Println(buildStringFromGroups(match.Groups, format))
		} else {
			fmt.Println(match)
		}
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
