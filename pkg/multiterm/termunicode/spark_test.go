package termunicode

import (
	"rare/pkg/testutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSparkUnicode(t *testing.T) {
	var sb strings.Builder
	SparkWrite(&sb, 0.0)
	SparkWrite(&sb, 0.1)
	SparkWrite(&sb, 0.2)
	SparkWrite(&sb, 0.3)
	SparkWrite(&sb, 0.4)
	SparkWrite(&sb, 0.5)
	SparkWrite(&sb, 0.6)
	SparkWrite(&sb, 0.7)
	SparkWrite(&sb, 0.8)
	SparkWrite(&sb, 0.9)
	SparkWrite(&sb, 1.0)

	assert.Equal(t, "__▁▂▃▄▄▅▆▇█", sb.String())
}

func TestSparkAscii(t *testing.T) {
	defer testutil.RestoreGlobals()
	testutil.SwitchGlobal(&UnicodeEnabled, false)

	var sb strings.Builder
	SparkWrite(&sb, 0.0)
	SparkWrite(&sb, 0.1)
	SparkWrite(&sb, 0.2)
	SparkWrite(&sb, 0.3)
	SparkWrite(&sb, 0.4)
	SparkWrite(&sb, 0.5)
	SparkWrite(&sb, 0.6)
	SparkWrite(&sb, 0.7)
	SparkWrite(&sb, 0.8)
	SparkWrite(&sb, 0.9)
	SparkWrite(&sb, 1.0)

	assert.Equal(t, "____...---^", sb.String())
}
