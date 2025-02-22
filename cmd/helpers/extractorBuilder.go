package helpers

import (
	"errors"
	"os"
	"rare/pkg/expressions"
	"rare/pkg/extractor"
	"rare/pkg/extractor/batchers"
	"rare/pkg/extractor/dirwalk"
	"rare/pkg/logger"
	"rare/pkg/matchers"
	"rare/pkg/matchers/dissect"
	"rare/pkg/matchers/fastregex"
	"runtime"
	"strings"

	"github.com/urfave/cli/v2"
)

const DefaultArgumentDescriptor = "<-|filename|glob...>"

const (
	cliCategoryRead     = "Input"
	cliCategoryOutput   = "Output"
	cliCategoryMatching = "Matching"
	cliCategoryTweaking = "Tweaking"
)

func BuildBatcherFromArguments(c *cli.Context) *batchers.Batcher {
	var (
		follow            = c.Bool("follow") || c.Bool("reopen")
		followTail        = c.Bool("tail")
		followReopen      = c.Bool("reopen")
		followPoll        = c.Bool("poll")
		concurrentReaders = c.Int("readers")
		gunzip            = c.Bool("gunzip")
		batchSize         = c.Int("batch")
		batchBuffer       = c.Int("batch-buffer")
		recursive         = c.Bool("recursive")
	)

	if batchSize < 1 {
		logger.Fatalf(ExitCodeInvalidUsage, "Batch size must be >= 1, is %d", batchSize)
	}
	if concurrentReaders < 1 {
		logger.Fatalf(ExitCodeInvalidUsage, "Must have at least 1 reader")
	}
	if followPoll && !follow {
		logger.Fatalf(ExitCodeInvalidUsage, "Follow (-f) must be enabled for --poll")
	}
	if followTail && !follow {
		logger.Fatalf(ExitCodeInvalidUsage, "Follow (-f) must be enabled for --tail")
	}

	fileglobs := c.Args().Slice()

	if len(fileglobs) == 0 || fileglobs[0] == "-" { // Read from stdin
		if gunzip {
			logger.Fatalln(ExitCodeInvalidUsage, "Cannot decompress (-z) with stdin")
		}
		if follow {
			logger.Println("Cannot follow a stdin stream, not a file")
		}
		return batchers.OpenReaderToChan("<stdin>", os.Stdin, batchSize, batchBuffer)
	} else if follow { // Read from source file
		if gunzip {
			logger.Println("Cannot combine -f and -z")
		}
		return batchers.TailFilesToChan(dirwalk.GlobExpand(fileglobs, recursive), batchSize, batchBuffer, followReopen, followPoll, followTail)
	} else { // Read (no-follow) source file(s)
		return batchers.OpenFilesToChan(dirwalk.GlobExpand(fileglobs, recursive), gunzip, concurrentReaders, batchSize, batchBuffer)
	}
}

func BuildExtractorFromArguments(c *cli.Context, batcher *batchers.Batcher) *extractor.Extractor {
	return BuildExtractorFromArgumentsEx(c, batcher, expressions.ArraySeparatorString, false)
}

func BuildExtractorFromArgumentsEx(c *cli.Context, batcher *batchers.Batcher, sep string, fullMatch bool) *extractor.Extractor {
	config := extractor.Config{
		Extract:   strings.Join(c.StringSlice("extract"), sep),
		Workers:   c.Int("workers"),
		FullMatch: fullMatch,
	}

	matcher, err := BuildMatcherFromArguments(c)
	if err != nil {
		logger.Fatalln(ExitCodeInvalidUsage, err)
	}
	config.Matcher = matcher

	ignoreSlice := c.StringSlice("ignore")
	if len(ignoreSlice) > 0 {
		ignoreExp, err := extractor.NewIgnoreExpressions(ignoreSlice...)
		if err != nil {
			logger.Fatalln(ExitCodeInvalidUsage, err)
		}
		config.Ignore = ignoreExp
	}

	ret, err := extractor.New(batcher.BatchChan(), &config)
	if err != nil {
		logger.Fatalln(ExitCodeInvalidUsage, err)
	}
	return ret
}

func BuildMatcherFromArguments(c *cli.Context) (matchers.Factory, error) {
	var (
		matchExpr   = c.String("match")
		dissectExpr = c.String("dissect")
		posix       = c.Bool("posix")
		ignoreCase  = c.Bool("ignore-case")
	)

	switch {
	case c.IsSet("match") && c.IsSet("dissect"):
		return nil, errors.New("match and dissect conflict")
	case c.IsSet("dissect"):
		d, err := dissect.CompileEx(dissectExpr, ignoreCase)
		if err != nil {
			return nil, err
		}
		return matchers.ToFactory(d), nil
	case c.IsSet("match"):
		if ignoreCase {
			matchExpr = "(?i)" + matchExpr
		}

		r, err := fastregex.CompileEx(matchExpr, posix)
		if err != nil {
			return nil, err
		}
		return matchers.ToFactory(r), nil
	default:
		return &matchers.AlwaysMatch{}, nil
	}
}

func getExtractorFlags() []cli.Flag {
	workerCount := runtime.NumCPU()/2 + 1

	return []cli.Flag{
		&cli.BoolFlag{
			Name:     "follow",
			Aliases:  []string{"f"},
			Category: cliCategoryRead,
			Usage:    "Read appended data as file grows",
		},
		&cli.BoolFlag{
			Name:     "reopen",
			Aliases:  []string{"F"},
			Category: cliCategoryRead,
			Usage:    "Same as -f, but will reopen recreated files",
		},
		&cli.BoolFlag{
			Name:     "poll",
			Category: cliCategoryRead,
			Usage:    "When following a file, poll for changes rather than using inotify",
		},
		&cli.BoolFlag{
			Name:     "tail",
			Aliases:  []string{"t"},
			Category: cliCategoryRead,
			Usage:    "When following a file, navigate to the end of the file to skip existing content",
		},
		&cli.BoolFlag{
			Name:     "gunzip",
			Aliases:  []string{"z"},
			Category: cliCategoryRead,
			Usage:    "Attempt to decompress file when reading",
		},
		&cli.BoolFlag{
			Name:     "recursive",
			Aliases:  []string{"R"},
			Category: cliCategoryRead,
			Usage:    "Recursively walk a non-globbing path and search for plain-files",
		},
		&cli.BoolFlag{
			Name:     "posix",
			Aliases:  []string{"p"},
			Category: cliCategoryMatching,
			Usage:    "Compile regex as against posix standard",
		},
		&cli.StringFlag{
			Name:     "match",
			Aliases:  []string{"m"},
			Category: cliCategoryMatching,
			Usage:    "Regex to create match groups to summarize on",
			Value:    ".*",
		},
		&cli.StringFlag{
			Name:     "dissect",
			Aliases:  []string{"d"},
			Category: cliCategoryMatching,
			Usage:    "Dissect expression create match groups to summarize on",
		},
		&cli.StringSliceFlag{
			Name:     "extract",
			Aliases:  []string{"e"},
			Category: cliCategoryMatching,
			Usage:    "Expression that will generate the key to group by. Specify multiple times for multi-dimensions or use {$} helper",
			Value:    cli.NewStringSlice("{0}"),
		},
		&cli.StringSliceFlag{
			Name:     "ignore",
			Aliases:  []string{"i"},
			Category: cliCategoryMatching,
			Usage:    "Ignore a match given a truthy expression (Can have multiple)",
		},
		&cli.BoolFlag{
			Name:     "ignore-case",
			Aliases:  []string{"I"},
			Category: cliCategoryMatching,
			Usage:    "Augment matcher to be case insensitive",
		},
		&cli.IntFlag{
			Name:     "batch",
			Category: cliCategoryTweaking,
			Usage:    "Specifies io batching size. Set to 1 for immediate input",
			Value:    1000,
		},
		&cli.IntFlag{
			Name:     "batch-buffer",
			Category: cliCategoryTweaking,
			Usage:    "Specifies how many batches to read-ahead. Impacts memory usage, can improve performance",
			Value:    workerCount * 2, // Keep 2 batches ready for each worker
		},
		&cli.IntFlag{
			Name:     "workers",
			Aliases:  []string{"w"},
			Category: cliCategoryTweaking,
			Usage:    "Set number of data processors",
			Value:    workerCount,
		},
		&cli.IntFlag{
			Name:     "readers",
			Aliases:  []string{"wr"},
			Category: cliCategoryTweaking,
			Usage:    "Sets the number of concurrent readers (Infinite when -f)",
			Value:    3,
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
