package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Exercise common commands for panics
func TestHistogram(t *testing.T) {
	testCommandSet(t, histogramCommand(),
		`-m (\d+) testdata/log.txt`,
		`-m (\d+) testdata/graph.txt`,
		`-m (\d+) --all testdata/graph.txt`,
		`-m (\d+) --scale log10 testdata/graph.txt`,
		`-o - -m (\d+) testdata/graph.txt`,
		`-z -m (\d+) testdata/log.txt.gz`,
	)
}

func TestHistogramRender(t *testing.T) {
	out, eout, err := testCommandCapture(histogramCommand(), `--snapshot -m "(\d+)" -e "{bucket {1} 10}" testdata/log.txt`)
	assert.NoError(t, err)
	assert.Equal(t, out, "0                   2         \n20                  1         \n\n\n\nMatched: 3 / 6 (Groups: 2)\n96 B (0 B/s) \n")
	assert.Equal(t, "", eout)
}

func TestHistogramCSV(t *testing.T) {
	out, eout, err := testCommandCapture(histogramCommand(), `-o - -m "(\d+)" -e "{bucket {1} 10}" testdata/log.txt`)
	assert.NoError(t, err)
	assert.Equal(t, "group,value\n0,2\n20,1\n", out)
	assert.Equal(t, "", eout)
}
