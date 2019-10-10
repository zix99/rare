package main

import (
	"fmt"
	"log"
	"os"

	"rare/cmd"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Version = fmt.Sprintf("%s, %s", version, buildSha)

	app.Commands = []cli.Command{
		*cmd.FilterCommand(),
		*cmd.HistogramCommand(),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
