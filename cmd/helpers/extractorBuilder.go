package helpers

import (
	"os"
	"rare/pkg/extractor"
	"runtime"

	"github.com/hpcloud/tail"
	"github.com/urfave/cli"
)

const DefaultArgumentDescriptor = "<-|filename|glob...>"

func BuildExtractorFromArguments(c *cli.Context) *extractor.Extractor {
	follow := c.Bool("follow") || c.Bool("reopen")
	followReopen := c.Bool("reopen")
	followPoll := c.Bool("poll")
	concurrentReaders := c.Int("readers")
	gunzip := c.Bool("gunzip")
	batchSize := c.Int("batch")
	recursive := c.Bool("recursive")
	config := extractor.Config{
		Posix:   c.Bool("posix"),
		Regex:   c.String("match"),
		Extract: c.String("extract"),
		Workers: c.Int("workers"),
	}

	ignoreSlice := c.StringSlice("ignore")
	if ignoreSlice != nil && len(ignoreSlice) > 0 {
		ignoreExp, err := extractor.NewIgnoreExpressions(ignoreSlice...)
		if err != nil {
			ErrLog.Panicln(err)
		}
		config.Ignore = ignoreExp
	}

	if batchSize < 1 {
		ErrLog.Fatalf("Batch size must be >= 1, is %d\n", batchSize)
	}

	fileglobs := c.Args()

	if len(fileglobs) == 0 || fileglobs[0] == "-" { // Read from stdin
		ret, err := extractor.New(extractor.ConvertReaderToStringChan("stdin", os.Stdin, batchSize), &config)
		if err != nil {
			ErrLog.Panicln(err)
		}
		StartFileReading("<stdin>")
		return ret
	} else if follow { // Read from source file
		if gunzip {
			ErrLog.Println("Cannot combine -f and -z")
		}

		tailChannels := make([]<-chan extractor.InputBatch, 0)
		for filename := range globExpand(fileglobs, recursive) {
			tail, err := tail.TailFile(filename, tail.Config{Follow: true, ReOpen: followReopen, Poll: followPoll})

			if err != nil {
				ErrLog.Fatal("Unable to open file: ", err)
			}
			tailChannels = append(tailChannels, tailLineToChan(filename, tail.Lines, batchSize))
			StartFileReading(filename)
		}

		ret, err := extractor.New(extractor.CombineChannels(tailChannels...), &config)
		if err != nil {
			ErrLog.Panicln(err)
		}
		return ret
	} else { // Read (no-follow) source file(s)
		ret, err := extractor.New(openFilesToChan(globExpand(fileglobs, recursive), gunzip, concurrentReaders, batchSize), &config)
		if err != nil {
			ErrLog.Panicln(err)
		}
		return ret
	}
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
		cli.BoolFlag{
			Name:  "recursive,R",
			Usage: "Recursively walk a non-globbing path and search for plain-files",
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
		DisableAndFlushLogBuffer()
		if originalAfter != nil {
			return originalAfter(c)
		}
		return nil
	}

	return &command
}
