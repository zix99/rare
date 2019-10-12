package cmd

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"rare/pkg/extractor"
	"runtime"
	"sync"

	"github.com/hpcloud/tail"
	"github.com/urfave/cli"
)

var stderrLog = log.New(os.Stderr, "[Log] ", 0)

func tailLineToChan(lines chan *tail.Line) <-chan string {
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

func openFilesToChan(filenames []string, gunzip bool, concurrency int) <-chan string {
	out := make(chan string, 128)
	sema := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for _, filename := range filenames {
		wg.Add(1)
		go func(goFilename string) {
			sema <- struct{}{}
			var file io.ReadCloser
			file, err := openFileToReader(goFilename, gunzip)
			if err != nil {
				stderrLog.Printf("Error opening file %s: %v\n", goFilename, err)
				return
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			bigBuf := make([]byte, 512*1024)
			scanner.Buffer(bigBuf, len(bigBuf))

			for scanner.Scan() {
				out <- scanner.Text()
			}
			<-sema
			wg.Done()
		}(filename)
	}

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

func buildExtractorFromArguments(c *cli.Context) *extractor.Extractor {
	follow := c.Bool("follow") || c.Bool("reopen")
	followReopen := c.Bool("reopen")
	followPoll := c.Bool("poll")
	concurrentReaders := c.Int("readers")
	gunzip := c.Bool("gunzip")
	config := extractor.Config{
		Posix:   c.Bool("posix"),
		Regex:   c.String("match"),
		Extract: c.String("extract"),
		Workers: c.Int("workers"),
	}

	if c.NArg() == 0 || c.Args().First() == "-" { // Read from stdin
		return extractor.New(extractor.ConvertReaderToStringChan(os.Stdin), &config)
	} else if follow { // Read from source file
		if gunzip {
			stderrLog.Println("Cannot combine -f and -z")
		}

		tailChannels := make([]<-chan string, 0)
		for _, filename := range globExpand(c.Args()) {
			tail, err := tail.TailFile(filename, tail.Config{Follow: true, ReOpen: followReopen, Poll: followPoll})

			if err != nil {
				stderrLog.Fatal("Unable to open file: ", err)
			}
			tailChannels = append(tailChannels, tailLineToChan(tail.Lines))
		}

		return extractor.New(extractor.CombineChannels(tailChannels...), &config)
	} else { // Read (no-follow) source file(s)
		return extractor.New(openFilesToChan(globExpand(c.Args()), gunzip, concurrentReaders), &config)
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
	}
	return append(flags, additionalFlags...)
}

func writeExtractorSummary(extractor *extractor.Extractor) {
	fmt.Fprintf(os.Stderr, "Matched: %d / %d\n", extractor.MatchedLines(), extractor.ReadLines())
}
