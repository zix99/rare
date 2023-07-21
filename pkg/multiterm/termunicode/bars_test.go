package termunicode

import (
	"rare/pkg/color"
	"rare/pkg/multiterm/termscaler"
	"rare/pkg/testutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteBar(t *testing.T) {
	var sb strings.Builder
	BarWriteScaled(&sb, termscaler.ScalerLinear.Scale(1, 0, 8), 1)
	assert.Equal(t, string(barUnicode[1]), sb.String())

	sb.Reset()
	BarWriteScaled(&sb, termscaler.ScalerLinear.Scale(10, 0, 8), 1)
	assert.Equal(t, string(fullBlock), sb.String())

	sb.Reset()
	BarWriteScaled(&sb, termscaler.ScalerLinear.Scale(5, 0, 8), 1)
	assert.Equal(t, string(barUnicode[5]), sb.String())
}

func TestWriteBarFallbacks(t *testing.T) {
	testutil.SwitchGlobal(&UnicodeEnabled, false)
	defer testutil.RestoreGlobals()

	var sb strings.Builder
	BarWriteScaled(&sb, termscaler.ScalerLinear.Scale(5, 0, 10), 10)
	assert.Equal(t, "|||||", sb.String())
}

func TestWriteBarStacked(t *testing.T) {
	defer testutil.RestoreGlobals()
	testutil.SwitchGlobal(&UnicodeEnabled, false)
	testutil.SwitchGlobal(&color.Enabled, false)

	var sb strings.Builder
	testutil.SwitchGlobal(&UnicodeEnabled, false)
	testutil.SwitchGlobal(&color.Enabled, false)

	BarWriteStacked(&sb, 10, 10, 1, 3, 2)
	assert.Equal(t, "011122", sb.String())

	testutil.SwitchGlobal(&color.Enabled, true)

	sb.Reset()
	BarWriteStacked(&sb, 10, 10, 1, 3, 2)
	assert.Equal(t, "\x1b[31m|\x1b[0m\x1b[32m|||\x1b[0m\x1b[33m||\x1b[0m", sb.String())

	testutil.SwitchGlobal(&UnicodeEnabled, true)
	sb.Reset()
	BarWriteStacked(&sb, 10, 10, 1, 3, 2)
	assert.Equal(t, "\x1b[31m█\x1b[0m\x1b[32m███\x1b[0m\x1b[33m██\x1b[0m", sb.String())
}

func TestBarWriteScaled(t *testing.T) {
	defer testutil.RestoreGlobals()
	testutil.SwitchGlobal(&UnicodeEnabled, false)
	var sb strings.Builder

	BarWriteScaled(&sb, 0.0, 6)
	assert.Equal(t, "", sb.String())
	sb.Reset()

	BarWriteScaled(&sb, 0.5, 6)
	assert.Equal(t, "|||", sb.String())
	sb.Reset()

	BarWriteScaled(&sb, 1.0, 6)
	assert.Equal(t, "||||||", sb.String())
	sb.Reset()

	testutil.SwitchGlobal(&UnicodeEnabled, true)
	BarWriteScaled(&sb, 0.45, 16)
	assert.Equal(t, "███████▏", sb.String())
	sb.Reset()

	BarWriteScaled(&sb, 1.0, 16)
	assert.Equal(t, "████████████████", sb.String())
	sb.Reset()
}

func TestBarKeyChar(t *testing.T) {
	defer testutil.RestoreGlobals()

	testutil.SwitchGlobal(&color.Enabled, false)
	testutil.SwitchGlobal(&UnicodeEnabled, false)
	assert.Equal(t, "0", BarKey(0))

	testutil.SwitchGlobal(&color.Enabled, true)
	testutil.SwitchGlobal(&UnicodeEnabled, false)
	assert.Equal(t, "\x1b[31m|\x1b[0m", BarKey(0))

	testutil.SwitchGlobal(&color.Enabled, false)
	testutil.SwitchGlobal(&UnicodeEnabled, true)
	assert.Equal(t, "0", BarKey(0))

	testutil.SwitchGlobal(&color.Enabled, true)
	testutil.SwitchGlobal(&UnicodeEnabled, true)
	assert.Equal(t, "\x1b[31m█\x1b[0m", BarKey(0))
}
