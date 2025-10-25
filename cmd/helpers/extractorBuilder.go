package helpers

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"slices"
	"strings"

	"github.com/zix99/rare/pkg/expressions"
	"github.com/zix99/rare/pkg/extractor"
	"github.com/zix99/rare/pkg/extractor/batchers"
	"github.com/zix99/rare/pkg/extractor/dirwalk"
	"github.com/zix99/rare/pkg/extractor/dirwalk/pathmatch"
	"github.com/zix99/rare/pkg/logger"
	"github.com/zix99/rare/pkg/matchers"
	"github.com/zix99/rare/pkg/matchers/dissect"
	"github.com/zix99/rare/pkg/matchers/fastregex"

	"github.com/urfave/cli/v2"
)

const DefaultArgumentDescriptor = "<-|filename|glob...>"

const (
	cliCategoryPath     = "Path"
	cliCategoryRead     = "Input"
	cliCategoryOutput   = "Output"
	cliCategoryMatching = "Matching"
	cliCategoryTweaking = "Tweaking"
)

func BuildBatcherFromArguments(c *cli.Context) *batchers.Batcher {
	batcher, _ := BuildBatcherFromArgumentsEx(c, c.Args().Slice()...)
	return batcher
}

func BuildBatcherFromArgumentsEx(c *cli.Context, fileglobs ...string) (*batchers.Batcher, dirwalk.Metrics) {
	var (
		follow            = c.Bool("follow") || c.Bool("reopen")
		followTail        = c.Bool("tail")
		followReopen      = c.Bool("reopen")
		followPoll        = c.Bool("poll")
		concurrentReaders = c.Int("readers")
		readersBuffer     = c.Int("readers-buffer")
		gunzip            = c.Bool("gunzip")
		batchSize         = c.Int("batch")
		batchBuffer       = c.Int("batch-buffer")
	)

	if batchSize < 1 {
		logger.Fatalf(ExitCodeInvalidUsage, "Batch size must be >= 1, is %d", batchSize)
	}
	if concurrentReaders < 1 {
		logger.Fatalf(ExitCodeInvalidUsage, "Must have at least 1 reader")
	}
	if readersBuffer < 1 {
		logger.Fatalf(ExitCodeInvalidUsage, "Reader buffer must be at least 1 byte")
	}
	if followPoll && !follow {
		logger.Fatalf(ExitCodeInvalidUsage, "Follow (-f) must be enabled for --poll")
	}
	if followTail && !follow {
		logger.Fatalf(ExitCodeInvalidUsage, "Follow (-f) must be enabled for --tail")
	}

	if len(fileglobs) == 0 || fileglobs[0] == "-" { // Read from stdin
		if gunzip {
			logger.Fatalln(ExitCodeInvalidUsage, "Cannot decompress (-z) with stdin")
		}
		if follow {
			logger.Println("Cannot follow a stdin stream, not a file")
		}
		return batchers.OpenReaderToChan("<stdin>", os.Stdin, batchSize, batchBuffer, readersBuffer), nil
	} else if follow { // Read from source file
		if gunzip {
			logger.Println("Cannot combine -f and -z")
		}
		walker := BuildPathWalkerFromArguments(c)
		return batchers.TailFilesToChan(walker.Walk(fileglobs...), batchSize, batchBuffer, readersBuffer, followReopen, followPoll, followTail), walker
	} else { // Read (no-follow) source file(s)
		walker := BuildPathWalkerFromArguments(c)
		return batchers.OpenFilesToChan(walker.Walk(fileglobs...), gunzip, concurrentReaders, batchSize, batchBuffer, readersBuffer), walker
	}
}

func BuildPathWalkerFromArguments(c *cli.Context) *dirwalk.Walker {
	var (
		include    = c.StringSlice("include")
		exclude    = c.StringSlice("exclude")
		excludeDir = c.StringSlice("exclude-dir")
	)

	includeSet, err := pathmatch.NewMatchSet(include...)
	if err != nil {
		logger.Fatal(ExitCodeInvalidUsage, err)
	}

	excludeSet, err := pathmatch.NewMatchSet(exclude...)
	if err != nil {
		logger.Fatal(ExitCodeInvalidUsage, err)
	}

	excludeDirSet, err := pathmatch.NewMatchSet(excludeDir...)
	if err != nil {
		logger.Fatal(ExitCodeInvalidUsage, err)
	}

	return &dirwalk.Walker{
		Filters: pathmatch.PathMatcher{
			Include:    includeSet,
			Exclude:    excludeSet,
			ExcludeDir: excludeDirSet,
		},
		Recursive:       c.Bool("recursive"),
		FollowSymLinks:  c.Bool("follow-symlinks"),
		ListSymLinks:    c.Bool("read-symlinks"),
		NoMountTraverse: c.Bool("mount"),
		OnTraverseError: func(err error) {
			logger.Print(err)
		},
	}
}

func BuildExtractorFromArguments(c *cli.Context, batcher *batchers.Batcher) *extractor.Extractor {
	return BuildExtractorFromArgumentsEx(c, batcher, expressions.ArraySeparatorString)
}

func BuildExtractorFromArgumentsEx(c *cli.Context, batcher *batchers.Batcher, sep string) *extractor.Extractor {
	config := extractor.Config{
		Extract: strings.Join(c.StringSlice("extract"), sep),
		Workers: c.Int("workers"),
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
		return matchers.NoFactory(d), nil
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
			EnvVars:  []string{"RARE_BATCH"},
		},
		&cli.IntFlag{
			Name:        "batch-buffer",
			Category:    cliCategoryTweaking,
			Usage:       "Specifies how many batches to read-ahead. Impacts memory usage, can improve performance",
			Value:       workerCount * 2, // Keep 2 batches ready for each worker
			DefaultText: fmt.Sprintf("workers*2 = %d", workerCount*2),
			EnvVars:     []string{"RARE_BATCH_BUFFER"},
		},
		&cli.IntFlag{
			Name:        "workers",
			Aliases:     []string{"w"},
			Category:    cliCategoryTweaking,
			Usage:       "Set number of data processors",
			Value:       workerCount,
			DefaultText: fmt.Sprintf("NumCPU/2+1 = %d", workerCount),
			EnvVars:     []string{"RARE_WORKERS"},
		},
		&cli.IntFlag{
			Name:     "readers",
			Aliases:  []string{"wr"},
			Category: cliCategoryTweaking,
			Usage:    "Sets the number of concurrent readers (Infinite when -f)",
			Value:    workerCount,
			EnvVars:  []string{"RARE_READERS"},
		},
		&cli.IntFlag{
			Name:     "readers-buffer",
			Aliases:  []string{"wrb"},
			Category: cliCategoryTweaking,
			Usage:    "How many bytes will be buffered per reader",
			Value:    1024 * 1024,
			EnvVars:  []string{"RARE_READERS_BUFFER"},
		},
	}
}

func GetWalkerFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:     "recursive",
			Aliases:  []string{"R"},
			Category: cliCategoryPath,
			Usage:    "Recursively walk a non-globbing path and search for plain-files",
		},
		&cli.StringSliceFlag{
			Name:     "include",
			Category: cliCategoryPath,
			Usage:    "Glob file patterns to include (eg. *.txt)",
			EnvVars:  []string{"RARE_INCLUDE"},
		},
		&cli.StringSliceFlag{
			Name:     "exclude",
			Category: cliCategoryPath,
			Usage:    "Glob file patterns to exclude (eg. *.txt)",
			EnvVars:  []string{"RARE_EXCLUDE"},
		},
		&cli.StringSliceFlag{
			Name:     "exclude-dir",
			Category: cliCategoryPath,
			Usage:    "Glob file patterns to exclude directories",
			EnvVars:  []string{"RARE_EXCLUDE_DIR"},
		},
		&cli.BoolFlag{
			Name:     "follow-symlinks",
			Aliases:  []string{"L"},
			Category: cliCategoryPath,
			Usage:    "Follow symbolic directory links",
			EnvVars:  []string{"RARE_FOLLOW_SYMLINKS"},
		},
		&cli.BoolFlag{
			Name:     "read-symlinks",
			Category: cliCategoryPath,
			Usage:    "Read files that are symbolic links",
			Value:    true,
		},
		&cli.BoolFlag{
			Name:     "mount",
			Category: cliCategoryPath,
			Usage:    "Don't descend directories on other filesystems (unix only)",
			Hidden:   !dirwalk.FeatureMountTraversal,
		},
	}
}

func AdaptCommandForExtractor(command cli.Command) *cli.Command {
	command.Flags = append(getExtractorFlags(), command.Flags...)
	command.Flags = append(GetWalkerFlags(), command.Flags...)
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

func ModifyArgOrPanic[T cli.Flag](cmd *cli.Command, name string, modifier func(T)) {
	for _, flag := range cmd.Flags {
		if slices.Contains(flag.Names(), name) {
			modifier(flag.(T))
			return
		}
	}
	panic("no flag change")
}
