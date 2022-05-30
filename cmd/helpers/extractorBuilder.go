package helpers

import (
	"os"
	"rare/pkg/expressions"
	"rare/pkg/extractor"
	"rare/pkg/extractor/batchers"
	"rare/pkg/extractor/dirwalk"
	"rare/pkg/logger"
	"runtime"
	"strings"

	"github.com/urfave/cli"
)

const DefaultArgumentDescriptor = "<-|filename|glob...>"

func BuildBatcherFromArguments(c *cli.Context) *batchers.Batcher {
	var (
		follow            = c.Bool("follow") || c.Bool("reopen")
		followReopen      = c.Bool("reopen")
		followPoll        = c.Bool("poll")
		concurrentReaders = c.Int("readers")
		gunzip            = c.Bool("gunzip")
		batchSize         = c.Int("batch")
		recursive         = c.Bool("recursive")
	)

	if batchSize < 1 {
		logger.Fatalf("Batch size must be >= 1, is %d", batchSize)
	}
	if concurrentReaders < 1 {
		logger.Fatalf("Must have at least 1 reader")
	}
	if followPoll && !follow {
		logger.Fatalf("Follow (-f) must be enabled for --poll")
	}

	fileglobs := c.Args()

	if len(fileglobs) == 0 || fileglobs[0] == "-" { // Read from stdin
		if gunzip {
			logger.Fatalln("Cannot decompress (-z) with stdin")
		}
		if follow {
			logger.Println("Cannot follow a stdin stream, not a file")
		}
		return batchers.OpenReaderToChan("<stdin>", os.Stdin, batchSize)
	} else if follow { // Read from source file
		if gunzip {
			logger.Println("Cannot combine -f and -z")
		}
		return batchers.TailFilesToChan(dirwalk.GlobExpand(fileglobs, recursive), batchSize, followReopen, followPoll)
	} else { // Read (no-follow) source file(s)
		return batchers.OpenFilesToChan(dirwalk.GlobExpand(fileglobs, recursive), gunzip, concurrentReaders, batchSize)
	}
}

func BuildExtractorFromArguments(c *cli.Context, batcher *batchers.Batcher) *extractor.Extractor {
	return BuildExtractorFromArgumentsEx(c, batcher, expressions.ArraySeparatorString)
}

func BuildExtractorFromArgumentsEx(c *cli.Context, batcher *batchers.Batcher, sep string) *extractor.Extractor {
	config := extractor.Config{
		Posix:   c.Bool("posix"),
		Regex:   c.String("match"),
		Extract: strings.Join(c.StringSlice("extract"), sep),
		Workers: c.Int("workers"),
	}

	if c.Bool("ignore-case") {
		config.Regex = "(?i)" + config.Regex
	}

	ignoreSlice := c.StringSlice("ignore")
	if len(ignoreSlice) > 0 {
		ignoreExp, err := extractor.NewIgnoreExpressions(ignoreSlice...)
		if err != nil {
			logger.Fatalln(err)
		}
		config.Ignore = ignoreExp
	}

	ret, err := extractor.New(batcher.BatchChan(), &config)
	if err != nil {
		logger.Fatalln(err)
	}
	return ret
}

func getExtractorFlags() []cli.Flag {
	return []cli.Flag{
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
		cli.StringSliceFlag{
			Name:  "extract,e",
			Usage: "Expression that will generate the key to group by. Specify multiple times for multi-dimensions or use {$} helper",
			Value: &cli.StringSlice{"{0}"},
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
		cli.BoolFlag{
			Name:  "recursive,R",
			Usage: "Recursively walk a non-globbing path and search for plain-files",
		},
		cli.BoolFlag{
			Name:  "ignore-case,I",
			Usage: "Augment regex to be case insensitive",
		},
	}
}

func AdaptCommandForExtractor(command cli.Command) *cli.Command {
	command.Flags = append(getExtractorFlags(), command.Flags...)
	if command.ArgsUsage == "" {
		command.ArgsUsage = DefaultArgumentDescriptor
	}

	// While this doesn't own the log, this is the last place
	// that has the option to flush the log buffer to sderr
	originalAfter := command.After
	command.After = func(c *cli.Context) error {
		logger.ImmediateLogs()
		if originalAfter != nil {
			return originalAfter(c)
		}
		return nil
	}

	return &command
}
