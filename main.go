package main

import (
	"log"
	"os"

	"reblurb/cmd"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Version = version

	app.Commands = []cli.Command{
		*cmd.FilterCommand(),
		*cmd.HistogramCommand(),
	}

	app.Flags = []cli.Flag{
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

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
