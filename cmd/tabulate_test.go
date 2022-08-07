package cmd

import (
	"testing"
)

func TestTabulate(t *testing.T) {
	testCommandSet(t, tabulateCommand(),
		`-m "(.+) (\d+)" -e "{$ {1} {2}}" testdata/graph.txt`,
		`-m "(.+) (\d+)" -e "{$ {1} {2}}" -x testdata/graph.txt`,
	)
}
