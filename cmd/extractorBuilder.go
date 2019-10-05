package cmd

import (
	"fmt"
	"log"
	"os"
	"rare/pkg/extractor"

	"github.com/hpcloud/tail"
	"github.com/urfave/cli"
)

func tailLineToString(lines chan *tail.Line) chan string {
	output := make(chan string)
	go func() {
		for {
			line := <-lines
			if line.Err != nil {
				break
			}
			output <- line.Text
		}
		close(output)
	}()
	return output
}

func buildExtractorFromArguments(c *cli.Context) *extractor.Extractor {
	follow := c.Bool("follow")
	config := extractor.Config{
		Posix:   c.Bool("posix"),
		Regex:   c.String("match"),
		Extract: c.String("extract"),
	}

	source := c.Args().First()

	if source == "" || source == "-" { // Read from stdin
		return extractor.NewExtractorReader(os.Stdin, &config)
	} else if follow { // Read from source file
		tail, err := tail.TailFile(source, tail.Config{Follow: true})
		if err != nil {
			log.Fatal("Unable to open file: ", err)
		}
		return extractor.NewExtractor(tailLineToString(tail.Lines), &config)
	} else { // Read (no-follow) source file(s)
		file, err := os.Open(source)
		if err != nil {
			log.Fatal(err)
		}

		return extractor.NewExtractorReader(file, &config)
	}
}

func buildExtractorFlags(additionalFlags ...cli.Flag) []cli.Flag {
	flags := []cli.Flag{
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
			Value: ".*",
		},
		cli.StringFlag{
			Name:  "extract,e",
			Usage: "Comparisons to extract",
		},
	}
	return append(flags, additionalFlags...)
}

func writeExtractorSummary(extractor *extractor.Extractor) {
	fmt.Fprintf(os.Stderr, "Matched: %d / %d\n", extractor.MatchedLines(), extractor.ReadLines())
}
