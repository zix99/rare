package cmd

import "github.com/urfave/cli"

func analyzeFunction(c *cli.Context) error {
	return nil
}

func AnalyzeCommand() *cli.Command {
	return &cli.Command{
		Name:      "analyze",
		Usage:     "Numerical analysis on a set of filtered data",
		Action:    analyzeFunction,
		ArgsUsage: "<-|filename|glob...>",
		Flags:     buildExtractorFlags(),
	}
}
