package cmd

import "testing"

func TestFilter(t *testing.T) {
	testCommandSet(t, filterCommand(),
		`-m \d+ testdata/log.txt`,
		`-m (\d+) testdata/log.txt`,
	)
}
