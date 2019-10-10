package cmd

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"rare/pkg/extractor"

	"github.com/hpcloud/tail"
	"github.com/urfave/cli"
)

var stderrLog = log.New(os.Stderr, "[Log] ", 0)

func tailLineToChan(lines chan *tail.Line) chan string {
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

func openFileToChan(filename string, gunzip bool) (chan string, error) {
	var file io.Reader
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	if gunzip {
		zfile, err := gzip.NewReader(file)
		if err != nil {
			stderrLog.Printf("Gunzip error: %v", err)
		} else {
			file = zfile
		}
	}

	return extractor.ConvertReaderToStringChan(file), nil
}

func buildExtractorFromArguments(c *cli.Context) *extractor.Extractor {
	follow := c.Bool("follow")
	gunzip := c.Bool("gunzip")
	config := extractor.Config{
		Posix:   c.Bool("posix"),
		Regex:   c.String("match"),
		Extract: c.String("extract"),
	}

	if c.NArg() == 0 || c.Args().First() == "-" { // Read from stdin
		return extractor.New(extractor.ConvertReaderToStringChan(os.Stdin), &config)
	} else if follow { // Read from source file
		tail, err := tail.TailFile(c.Args().First(), tail.Config{Follow: true})

		if err != nil {
			stderrLog.Fatal("Unable to open file: ", err)
		}
		if gunzip {
			stderrLog.Println("Cannot combine -f and -z")
		}
		return extractor.New(tailLineToChan(tail.Lines), &config)
	} else { // Read (no-follow) source file(s)
		fchan, err := openFileToChan(c.Args().First(), gunzip)
		if err != nil {
			stderrLog.Fatal(err)
		}

		return extractor.New(fchan, &config)
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
			Value: "{0}",
		},
		cli.BoolFlag{
			Name:  "gunzip,z",
			Usage: "Attempt to decompress file when reading",
		},
	}
	return append(flags, additionalFlags...)
}

func writeExtractorSummary(extractor *extractor.Extractor) {
	fmt.Fprintf(os.Stderr, "Matched: %d / %d\n", extractor.MatchedLines(), extractor.ReadLines())
}
