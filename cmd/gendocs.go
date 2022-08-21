//go:build !urfave_cli_no_docs

package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func gendocCommand() *cli.Command {
	return &cli.Command{
		Name:   "_gendoc",
		Hidden: true,
		Usage:  "Generates documentation",
		Action: func(c *cli.Context) error {
			var text string
			if c.Bool("man") {
				text, _ = c.App.ToMan()
			} else {
				text, _ = c.App.ToMarkdown()
			}
			fmt.Print(text)
			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "man",
				Usage: "manpage syntax",
			},
		},
	}
}

func init() {
	commands = append(commands, gendocCommand())
}
