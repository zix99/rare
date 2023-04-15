package termunicode

import (
	"rare/pkg/color"
	"rare/pkg/testutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteBar(t *testing.T) {
	var sb strings.Builder
	BarWrite(&sb, 1, 8, 1)
	assert.Equal(t, string(barUnicode[1]), sb.String())

	sb.Reset()
	BarWrite(&sb, 10, 8, 1)
	assert.Equal(t, string(fullBlock), sb.String())

	assert.Equal(t, string(barUnicode[5]), BarString(5, 8, 1))

	sb.Reset()
	BarWriteFull(&sb, 1, 8, 10)
	assert.Equal(t, string(fullBlock), sb.String())
}

func TestWriteBarFallbacks(t *testing.T) {
	UnicodeEnabled = false

	assert.Equal(t, "|||||", BarString(5, 10, 10))

	UnicodeEnabled = true
}

func TestWriteBarStacked(t *testing.T) {
	defer testutil.RevertGlobals()
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

func TestBarKeyChar(t *testing.T) {
	defer testutil.RevertGlobals()
	testutil.SwitchGlobal(&color.Enabled, true)
}
