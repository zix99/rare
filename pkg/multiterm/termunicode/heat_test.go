package termunicode

import (
	"rare/pkg/color"
	"rare/pkg/multiterm/termscaler"
	"rare/pkg/testutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeatLinear(t *testing.T) {
	defer testutil.RestoreGlobals()
	testutil.SwitchGlobal(&color.Enabled, true)

	var sb strings.Builder

	HeatWriteLinear(&sb, 0, 0, 0)
	HeatWriteLinear(&sb, -5, 0, 10)
	HeatWriteLinear(&sb, 15, 0, 10)
	HeatWriteLinear(&sb, 0, 0, 10)
	HeatWriteLinear(&sb, 5, 0, 10)
	HeatWriteLinear(&sb, 10, 0, 10)
	HeatWriteLinear(&sb, 5, 5, 2)
	HeatWriteLinear(&sb, 750, 500, 1000)

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

func TestHeatLog10(t *testing.T) {
	defer testutil.RestoreGlobals()
	testutil.SwitchGlobal(&color.Enabled, false)

	var sb strings.Builder

	HeatWrite(&sb, termscaler.ScalerLog10.Scale(0, 0, 10))
	HeatWrite(&sb, termscaler.ScalerLog10.Scale(2, 0, 10))
	HeatWrite(&sb, termscaler.ScalerLog10.Scale(5, 0, 10))
	HeatWrite(&sb, termscaler.ScalerLog10.Scale(10, 0, 10))

	assert.Equal(t, "-269", sb.String())
}
