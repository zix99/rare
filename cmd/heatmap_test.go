package cmd

import "testing"

func TestHeatmap(t *testing.T) {
	testCommandSet(t, heatmapCommand(),
		`-m "(.+) (\d+)" -e "{$ {1} {2}}" testdata/graph.txt`,
	)
}
