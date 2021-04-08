package cmd

import (
	"fmt"
	"os"
	"rare/docs"
	"rare/pkg/markdowncli"
	"strings"

	"github.com/urfave/cli"
)

func docsFunction(c *cli.Context) error {
	docname := strings.ToLower(c.Args().First())

	if docname == "" || docname == "list" {
		fmt.Println("Available Docs:")
		entries, _ := docs.DocFS.ReadDir(".")
		for _, entry := range entries {
			fmt.Printf("  %s\n", strings.Title(strings.TrimSuffix(entry.Name(), ".md")))
		}

	} else if file, err := docs.DocFS.Open(docname + ".md"); err == nil {
		if err != nil {
			fmt.Println(err)
		} else {
			markdowncli.WriteMarkdownToTerm(os.Stdout, file)
		}
	} else {
		fmt.Printf("Error: No such doc %s\n", docname)
	}
	return nil
}

func docsCommand() *cli.Command {
	return &cli.Command{
		Name:      "docs",
		Usage:     "Access help documentation",
		ArgsUsage: "[doc]",
		Action:    docsFunction,
	}
}
