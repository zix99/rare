package cmd

import "github.com/urfave/cli/v2"

var (
	cmdCatAnalyze   = "Analyze"
	cmdCatVisualize = "Visualize"
	cmdCatHelp      = "Help"
)

var commands []*cli.Command = []*cli.Command{
	filterCommand(),
	searchCommand(),
	histogramCommand(),
	heatmapCommand(),
	sparkCommand(),
	bargraphCommand(),
	analyzeCommand(),
	tabulateCommand(),
	reduceCommand(),
	docsCommand(),
	expressionCommand(),
	walkCommand(),
}

func GetSupportedCommands() []*cli.Command {
	return commands
}
