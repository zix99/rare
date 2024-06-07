package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeatmap(t *testing.T) {
	testCommandSet(t, heatmapCommand(),
		`-m "(.+) (\d+)" -e "{$ {1} {2}}" testdata/graph.txt`,
		`-o - -m "(.+) (\d+)" -e "{$ {1} {2}}" testdata/graph.txt`,
	)
}

func TestHeatmapLinear(t *testing.T) {
	out, eout, err := testCommandCapture(heatmapCommand(), `--snapshot -m "(.+) (.+)" -e "{$ {1} {2}}" testdata/heat.txt`)
	assert.NoError(t, err)
	assert.Empty(t, eout)
	assert.Contains(t, out, " - 0    2 1    4 2    6 3    9 4\n a..\nx -22\ny 224\nz 9--\nMatched: 10 / 10 (R: 3; C: 3)\n39 B (0 B/s)")
}

func TestHeatmapLog2(t *testing.T) {
	out, eout, err := testCommandCapture(heatmapCommand(), `--snapshot -m "(.+) (.+)" -e "{$ {1} {2}}" --scale log2 testdata/heat.txt`)
	assert.NoError(t, err)
	assert.Empty(t, eout)
	assert.Contains(t, out, " - 1    4 2    7 3    9 4\n a..\nx ---\ny --4\nz 9--\nMatched: 10 / 10 (R: 3; C: 3)\n39 B (0 B/s)")
}
