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
		*cmd.HistogramCommand(),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
