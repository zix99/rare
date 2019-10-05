package cmd

import (
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
	follow := c.GlobalBool("follow")
	config := extractor.Config{
		Posix:   c.GlobalBool("posix"),
		Regex:   c.GlobalString("match"),
		Extract: c.GlobalString("extract"),
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
