package termunicode

import (
	"rare/pkg/color"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeatLinear(t *testing.T) {
	var sb strings.Builder

	color.Enabled = true
	HeatWriteLinear(&sb, 0, 0, 0)
	HeatWriteLinear(&sb, -5, 0, 10)
	HeatWriteLinear(&sb, 15, 0, 10)
	HeatWriteLinear(&sb, 0, 0, 10)
	HeatWriteLinear(&sb, 5, 0, 10)
	HeatWriteLinear(&sb, 10, 0, 10)
	HeatWriteLinear(&sb, 5, 5, 2)
	HeatWriteLinear(&sb, 750, 500, 1000)
	color.Enabled = false

	assert.Equal(t, "\x1b[38;5;16m█\x1b[0m\x1b[38;5;16m█\x1b[0m\x1b[38;5;196m█\x1b[0m\x1b[38;5;16m█\x1b[0m\x1b[38;5;93m█\x1b[0m\x1b[38;5;196m█\x1b[0m\x1b[38;5;16m█\x1b[0m\x1b[38;5;93m█\x1b[0m", sb.String())
}

func TestHeatNoColor(t *testing.T) {
	var sb strings.Builder

	HeatWriteLinear(&sb, 0, 0, 0)
	HeatWriteLinear(&sb, -5, 0, 10)
	HeatWriteLinear(&sb, 15, 0, 10)
	HeatWriteLinear(&sb, 0, 0, 10)
	HeatWriteLinear(&sb, 6, 0, 10)
	HeatWriteLinear(&sb, 10, 0, 10)
	HeatWriteLinear(&sb, 5, 5, 2)
	HeatWriteLinear(&sb, 750, 500, 1000)
	HeatWriteLinear(&sb, 0, 0, 1)
	HeatWriteLinear(&sb, 1, 0, 1)

	assert.Equal(t, "--9-59-4-9", sb.String())
}
