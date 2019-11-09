package helpers

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"rare/pkg/color"
	"rare/pkg/extractor"
	"rare/pkg/humanize"
	"runtime"
	"sync"

	"github.com/hpcloud/tail"
	"github.com/urfave/cli"
)

const DefaultArgumentDescriptor = "<-|filename|glob...>"

var stderrLog = log.New(os.Stderr, "[Log] ", 0)

func tailLineToChan(lines chan *tail.Line) <-chan []string {
	output := make(chan []string)
	go func() {
		for {
			line := <-lines
			if line == nil || line.Err != nil {
				break
			}
			// Don't batch when tailing files
			output <- []string{line.Text}
		}
		close(output)
	}()
	return output
}

func openFileToReader(filename string, gunzip bool) (io.ReadCloser, error) {
	var file io.ReadCloser
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	if gunzip {
		zfile, err := gzip.NewReader(file)
		if err != nil {
			stderrLog.Printf("Gunzip error for file %s: %v\n", filename, err)
		} else {
			file = zfile
		}
	}

	return file, nil
}

func openFilesToChan(filenames []string, gunzip bool, concurrency int, batchSize int) <-chan []string {
	out := make(chan []string, 128)
	sema := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	wg.Add(len(filenames))
	IncSourceCount(len(filenames))

	// Load as many files as the sema allows
	go func() {
		for _, filename := range filenames {
			sema <- struct{}{}

			go func(goFilename string) {
				var file io.ReadCloser
				file, err := openFileToReader(goFilename, gunzip)
				if err != nil {
					stderrLog.Printf("Error opening file %s: %v\n", goFilename, err)
					return
				}
				defer file.Close()
				StartFileReading(goFilename)

				scanner := bufio.NewScanner(file)
				bigBuf := make([]byte, 512*1024)
				scanner.Buffer(bigBuf, len(bigBuf))

				batch := make([]string, 0, batchSize)
				for scanner.Scan() {
					batch = append(batch, scanner.Text())
					if len(batch) >= batchSize {
						out <- batch
						batch = make([]string, 0, batchSize)
					}
				}
				if len(batch) > 0 {
					out <- batch
				}

				<-sema
				wg.Done()
				StopFileReading(goFilename)
			}(filename)
		}
	}()

	// Wait on all files, and close chan
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func globExpand(paths []string) []string {
	out := make([]string, 0)
	for _, path := range paths {
		expanded, err := filepath.Glob(path)
		if err != nil {
			stderrLog.Printf("Path error: %v\n", err)
		} else {
			out = append(out, expanded...)
		}
	}
	return out
}

func BuildExtractorFromArguments(c *cli.Context) *extractor.Extractor {
	follow := c.Bool("follow") || c.Bool("reopen")
	followReopen := c.Bool("reopen")
	followPoll := c.Bool("poll")
	concurrentReaders := c.Int("readers")
	gunzip := c.Bool("gunzip")
	batchSize := c.Int("batch")
	config := extractor.Config{
		Posix:   c.Bool("posix"),
		Regex:   c.String("match"),
		Extract: c.String("extract"),
		Workers: c.Int("workers"),
	}

	ignoreSlice := c.StringSlice("ignore")
	if ignoreSlice != nil && len(ignoreSlice) > 0 {
		ignoreExp, err := extractor.NewIgnoreExpressions(ignoreSlice)
		if err != nil {
			log.Panicln(err)
		}
		config.Ignore = ignoreExp
	}

	if batchSize < 1 {
		stderrLog.Fatalf("Batch size must be >= 1, is %d\n", batchSize)
	}

	if c.NArg() == 0 || c.Args().First() == "-" { // Read from stdin
		ret, err := extractor.New(extractor.ConvertReaderToStringChan(os.Stdin), &config)
		if err != nil {
			log.Panicln(err)
		}
		StartFileReading("<stdin>")
		return ret
	} else if follow { // Read from source file
		if gunzip {
			stderrLog.Println("Cannot combine -f and -z")
		}

		tailChannels := make([]<-chan []string, 0)
		for _, filename := range globExpand(c.Args()) {
			tail, err := tail.TailFile(filename, tail.Config{Follow: true, ReOpen: followReopen, Poll: followPoll})

			if err != nil {
				stderrLog.Fatal("Unable to open file: ", err)
			}
			tailChannels = append(tailChannels, tailLineToChan(tail.Lines))
			StartFileReading(filename)
		}

		ret, err := extractor.New(extractor.CombineChannels(tailChannels...), &config)
		if err != nil {
			log.Panicln(err)
		}
		return ret
	} else { // Read (no-follow) source file(s)
		ret, err := extractor.New(openFilesToChan(globExpand(c.Args()), gunzip, concurrentReaders, batchSize), &config)
		if err != nil {
			log.Panicln(err)
		}
		return ret
	}
}

func BuildExtractorFlags(additionalFlags ...cli.Flag) []cli.Flag {
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
		cli.IntFlag{
			Name:  "batch",
			Usage: "Specifies io batching size. Set to 1 for immediate input",
			Value: 1000,
		},
		cli.IntFlag{
			Name:  "workers,w",
			Usage: "Set number of data processors",
			Value: runtime.NumCPU()/2 + 1,
		},
		cli.IntFlag{
			Name:  "readers,wr",
			Usage: "Sets the number of concurrent readers (Infinite when -f)",
			Value: 3,
		},
		cli.StringSliceFlag{
			Name:  "ignore,i",
			Usage: "Ignore a match given a truthy expression (Can have multiple)",
		},
	}
	return append(flags, additionalFlags...)
}

func WriteExtractorSummary(extractor *extractor.Extractor) {
	fmt.Fprintf(os.Stderr, "Matched: %s / %s",
		color.Wrapi(color.BrightGreen, humanize.Hi(extractor.MatchedLines())),
		color.Wrapi(color.BrightWhite, humanize.Hi(extractor.ReadLines())))
	if extractor.IgnoredLines() > 0 {
		fmt.Fprintf(os.Stderr, " (Ignored: %s)", color.Wrapi(color.Red, humanize.Hi(extractor.IgnoredLines())))
	}
	fmt.Fprintf(os.Stderr, "\n")
}
