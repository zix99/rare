package cmd

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

func buildStringFromGroups(matches []string, format string) string {
	var temp string = format
	for idx, val := range matches {
		temp = strings.Replace(temp, fmt.Sprintf("$%d", idx), val, -1)
	}
	return temp
}

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
			fmt.Println(buildStringFromGroups(match.Groups, format))
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
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "follow,f",
				Usage: "Read appended data as file grows",
			},
			cli.BoolFlag{
				Name:  "posix,p",
				Usage: "Compile regex as against posix standard",
			},
			cli.StringFlag{
				Name:  "match,m",
				Usage: "Regex to create match groups to summarize on",
			},
			cli.StringFlag{
				Name:  "filter",
				Usage: "Filters incoming lines without creating matches",
			},
			cli.StringFlag{
				Name:  "extract,e",
				Usage: "Comparisons to extract",
			},
		},
	}
}
