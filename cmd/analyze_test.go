package cmd

import (
	"rare/cmd/helpers"
	"testing"
)

func TestAnalyze(t *testing.T) {
	testCommandSet(t, analyzeCommand(),
		`-m (\d+) testdata/graph.txt`,
		`-x -m (\d+) testdata/graph.txt`,
	)
}

func TestAnalyzeParseFatals(t *testing.T) {
	catchLogFatal(t, helpers.ExitCodeInvalidUsage, func() {
		testCommand(analyzeCommand(), "--quantile bla testdata/graph.txt")
	})
}
