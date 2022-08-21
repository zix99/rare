package cmd

import "github.com/urfave/cli/v2"

var commands []*cli.Command = []*cli.Command{
	filterCommand(),
	histogramCommand(),
	heatmapCommand(),
	bargraphCommand(),
	analyzeCommand(),
	tabulateCommand(),
	docsCommand(),
}

func GetSupportedCommands() []*cli.Command {
	return commands
}
