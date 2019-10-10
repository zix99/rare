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
	follow := c.Bool("follow") || c.Bool("reopen")
	followReopen := c.Bool("reopen")
	followPoll := c.Bool("poll")
	gunzip := c.Bool("gunzip")
	config := extractor.Config{
		Posix:   c.Bool("posix"),
		Regex:   c.String("match"),
		Extract: c.String("extract"),
	}

	if c.NArg() == 0 || c.Args().First() == "-" { // Read from stdin
		return extractor.New(extractor.ConvertReaderToStringChan(os.Stdin), &config)
	} else if follow { // Read from source file
		if gunzip {
			stderrLog.Println("Cannot combine -f and -z")
		}

		tailChannels := make([]chan string, 0)
		for _, filename := range c.Args() {
			tail, err := tail.TailFile(filename, tail.Config{Follow: true, ReOpen: followReopen, Poll: followPoll})

			if err != nil {
				stderrLog.Fatal("Unable to open file: ", err)
			}
			tailChannels = append(tailChannels, tailLineToChan(tail.Lines))
		}

		return extractor.New(extractor.CombineChannels(tailChannels...), &config)
	} else { // Read (no-follow) source file(s)
		readChannels := make([]chan string, 0)
		for _, filename := range c.Args() {
			fchan, err := openFileToChan(filename, gunzip)
			if err != nil {
				stderrLog.Fatal(err)
			}
			readChannels = append(readChannels, fchan)
		}

		return extractor.New(extractor.CombineChannels(readChannels...), &config)
	}
}

func buildExtractorFlags(additionalFlags ...cli.Flag) []cli.Flag {
	flags := []cli.Flag{
		cli.BoolFlag{
			Name:  "follow,f",
			Usage: "Read appended data as file grows",
		},
		cli.BoolFlag{
			Name:  "reopen,F",
			Usage: "Same as -f, but will reopen recreated files",
		},
		cli.BoolFlag{
			Name:  "poll",
			Usage: "When following a file, poll for changes rather than using inotify",
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
			Usage: "Expression that will generate the key to group by",
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
