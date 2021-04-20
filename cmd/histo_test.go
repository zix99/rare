package cmd

import (
	"testing"
)

func TestHistogram(t *testing.T) {
	testCommandSet(t, histogramCommand(),
		`-m (\d+) testdata/log.txt`,
		`-m (\d+) testdata/graph.txt`,
		`-z -m (\d+) testdata/log.txt.gz`,
	)
}
