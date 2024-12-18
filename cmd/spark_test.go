package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSparkline(t *testing.T) {
	testCommandSet(t, sparkCommand(),
		`-m "(.+) (\d+)" -e "{$ {1} {2}}" testdata/graph.txt`,
		`-o - -m "(.+) (\d+)" -e "{$ {1} {2}}" testdata/graph.txt`,
	)
}

func TestSparklineWithTrim(t *testing.T) {
	out, eout, err := testCommandCapture(sparkCommand(), `--snapshot -m "(.+) (.+)" -e {1} -e {2} --cols 2 testdata/heat.txt`)

	assert.NoError(t, err)
	assert.Empty(t, eout)
	assert.Contains(t, out, "  First bc Last \ny 1     _█ 2    \nx 1     __ 1    \nMatched: 10 / 10 (R: 2; C: 2)")
}
