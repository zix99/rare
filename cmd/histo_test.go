package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHistogram(t *testing.T) {
	testCommandSet(t, histogramCommand(),
		`-m (\d+) testdata/log.txt`,
		`-m (\d+) testdata/graph.txt`,
		`-z -m (\d+) testdata/log.txt.gz`,
	)
}

func TestHistogramRender(t *testing.T) {
	out, eout, err := testCommandCapture(histogramCommand(), `-m "(\d+)" -e "{bucket {1} 10}" testdata/log.txt`)
	assert.NoError(t, err)
	assert.Contains(t, out, "Matched: 3 / 6 (Groups: 2)")
	assert.Equal(t, "", eout)
}
