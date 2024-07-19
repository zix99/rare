//go:build !urfave_cli_no_docs

package cmd

import (
	"fmt"
	"strings"

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
			fmt.Print(strings.ReplaceAll(text, "\x00", "")) //HACK: Some null characters are in generated docs (from array sep?)
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
