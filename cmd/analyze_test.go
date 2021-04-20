package cmd

import "testing"

func TestAnalyze(t *testing.T) {
	testCommandSet(t, analyzeCommand(),
		`-m (\d+) testdata/graph.txt`,
		`-x -m (\d+) testdata/graph.txt`,
	)
}
