package main

import (
	"log"
	"os"

	"rare/cmd"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Version = version

	app.Commands = []cli.Command{
		*cmd.FilterCommand(),
		*cmd.HistogramCommand(),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
