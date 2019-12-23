package cmd

import "github.com/urfave/cli"

func GetSupportedCommands() []cli.Command {
	return []cli.Command{
		*filterCommand(),
		*histogramCommand(),
		*analyzeCommand(),
		*tabulateCommand(),
		*docsCommand(),
		*distCommand(),
	}
}
