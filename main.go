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

	app.Usage = "A regex parser and extractor"

	app.Version = fmt.Sprintf("%s, %s", version, buildSha)

	app.Description = `Aggregate and display information parsed from text files using
	regex and a simple handlebars-like expression syntax.
	
	https://github.com/zix99/rare`

	app.Copyright = `rare  Copyright (C) 2019 Chris LaPointe
    This program comes with ABSOLUTELY NO WARRANTY.
    This is free software, and you are welcome to redistribute it
	under certain conditions`

	app.Commands = []cli.Command{
		*cmd.FilterCommand(),
		*cmd.HistogramCommand(),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
