package cmd

import "github.com/urfave/cli/v2"

var (
	cmdCatAnalyze   = "Analyze"
	cmdCatVisualize = "Visualize"
	cmdCatHelp      = "Help"
)

var commands []*cli.Command = []*cli.Command{
	filterCommand(),
	reduceCommand(),
	histogramCommand(),
	heatmapCommand(),
	bargraphCommand(),
	analyzeCommand(),
	tabulateCommand(),
	docsCommand(),
	expressionCommand(),
}

func GetSupportedCommands() []*cli.Command {
	return commands
}
