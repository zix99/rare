package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"rare/cmd/helpers"
	"rare/docs"
	"rare/pkg/markdowncli"
	"strings"

	"github.com/urfave/cli/v2"
)

func docsFunction(c *cli.Context) error {
	docname := strings.ToLower(c.Args().First())

	if docname == "" || docname == "list" {
		listDocFiles()
	} else if file, err := openDocFileByPartialName(docname); err == nil {
		var buf bytes.Buffer
		markdowncli.WriteMarkdownToBuf(&buf, file)
		if c.Bool("no-pager") || helpers.TryWritePager(&buf) != nil {
			io.Copy(os.Stdout, &buf)
		}
	} else {
		return cli.NewExitError(fmt.Sprintf("No such doc '%s'", docname), helpers.ExitCodeInvalidUsage)
	}

	return nil
}

func listDocFiles() {
	fmt.Println("Available Docs:")
	entries, _ := docs.DocFS.ReadDir(docs.BasePath)
	for _, entry := range entries {
		fmt.Printf("  %s\n", strings.TrimSuffix(entry.Name(), ".md"))
	}
}

func openDocFileByPartialName(s string) (fs.File, error) {
	s = strings.ToLower(s)

	// Try exact match
	if f, err := docs.DocFS.Open(docs.BasePath + "/" + s + ".md"); err == nil {
		return f, nil
	}

	// Search for prefix of name
	entries, _ := docs.DocFS.ReadDir(docs.BasePath)
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), s) {
			fullPath := docs.BasePath + "/" + entry.Name()
			f, err := docs.DocFS.Open(fullPath)
			if err == nil {
				return f, nil
			}
		}
	}

	return nil, errors.New("no such file")
}

func docsCommand() *cli.Command {
	return &cli.Command{
		Name:      "docs",
		Usage:     "Access detailed documentation",
		ArgsUsage: "[doc]",
		Action:    docsFunction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "no-pager",
				Aliases: []string{"n"},
				Usage:   "Don't use pager to view documentation",
			},
		},
	}
}
