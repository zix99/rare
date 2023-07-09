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

	HeatWrite(&sb, termscaler.ScalerLinear.Scale(0, 0, 0))
	HeatWrite(&sb, termscaler.ScalerLinear.Scale(-5, 0, 10))
	HeatWrite(&sb, termscaler.ScalerLinear.Scale(15, 0, 10))
	HeatWrite(&sb, termscaler.ScalerLinear.Scale(0, 0, 10))
	HeatWrite(&sb, termscaler.ScalerLinear.Scale(5, 0, 10))
	HeatWrite(&sb, termscaler.ScalerLinear.Scale(10, 0, 10))
	HeatWrite(&sb, termscaler.ScalerLinear.Scale(5, 5, 2))
	HeatWrite(&sb, termscaler.ScalerLinear.Scale(750, 500, 1000))

	assert.Equal(t, "\x1b[38;5;16m█\x1b[0m\x1b[38;5;16m█\x1b[0m\x1b[38;5;196m█\x1b[0m\x1b[38;5;16m█\x1b[0m\x1b[38;5;93m█\x1b[0m\x1b[38;5;196m█\x1b[0m\x1b[38;5;16m█\x1b[0m\x1b[38;5;93m█\x1b[0m", sb.String())
}

func TestHeatNoColor(t *testing.T) {
	var sb strings.Builder

	HeatWrite(&sb, termscaler.ScalerLinear.Scale(0, 0, 0))
	HeatWrite(&sb, termscaler.ScalerLinear.Scale(-5, 0, 10))
	HeatWrite(&sb, termscaler.ScalerLinear.Scale(15, 0, 10))
	HeatWrite(&sb, termscaler.ScalerLinear.Scale(0, 0, 10))
	HeatWrite(&sb, termscaler.ScalerLinear.Scale(6, 0, 10))
	HeatWrite(&sb, termscaler.ScalerLinear.Scale(10, 0, 10))
	HeatWrite(&sb, termscaler.ScalerLinear.Scale(5, 5, 2))
	HeatWrite(&sb, termscaler.ScalerLinear.Scale(750, 500, 1000))
	HeatWrite(&sb, termscaler.ScalerLinear.Scale(0, 0, 1))
	HeatWrite(&sb, termscaler.ScalerLinear.Scale(1, 0, 1))

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
