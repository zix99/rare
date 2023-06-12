package cmd

import "testing"

func TestBarGraph(t *testing.T) {
	testCommandSet(t, bargraphCommand(),
		`-m "(.+) (\d+)" -e "{$ {1} {2}}" testdata/graph.txt`,
		`-o - -m "(.+) (\d+)" -e "{$ {1} {2}}" testdata/graph.txt`,
	)
}
