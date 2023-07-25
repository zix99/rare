package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBarGraph(t *testing.T) {
	testCommandSet(t, bargraphCommand(),
		`-m "(.+) (\d+)" -e "{$ {1} {2}}" testdata/graph.txt`,
		`-o - -m "(.+) (\d+)" -e "{$ {1} {2}}" testdata/graph.txt`,
		`-o - -m "(.+) (\d+)" -e "{$ {1} {2}}" --scale log10 testdata/graph.txt`,
	)
}

func TestBarGraphCantScaleAndStack(t *testing.T) {
	err := testCommand(bargraphCommand(), "--stacked --scale log10 testdata/graph.txt")
	assert.Error(t, err)
}
