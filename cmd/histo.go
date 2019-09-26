package cmd

import (
	"fmt"

	"github.com/urfave/cli"
)

func histoFunction(c *cli.Context) error {
	fmt.Println("Howdy")
	return nil
}

func HistogramCommand() *cli.Command {
	return &cli.Command{
		Name:   "histo",
		Usage:  "Summarize results by extracting them to a histogram",
		Action: histoFunction,
	}
}
