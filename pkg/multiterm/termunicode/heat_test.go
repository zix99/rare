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

	assert.Equal(t, []byte{
		0x1b, 0x5b, 0x33, 0x38, 0x3b, 0x35, 0x3b, 0x31, 0x36, 0x6d, 0xe2, 0x96, 0x88, 0x1b, 0x5b, 0x30, 0x6d,
		0x1b, 0x5b, 0x33, 0x38, 0x3b, 0x35, 0x3b, 0x31, 0x36, 0x6d, 0xe2, 0x96, 0x88, 0x1b, 0x5b, 0x30, 0x6d,
		0x1b, 0x5b, 0x33, 0x38, 0x3b, 0x35, 0x3b, 0x31, 0x39, 0x37, 0x6d, 0xe2, 0x96, 0x88, 0x1b, 0x5b, 0x30, 0x6d,
		0x1b, 0x5b, 0x33, 0x38, 0x3b, 0x35, 0x3b, 0x31, 0x36, 0x6d, 0xe2, 0x96, 0x88, 0x1b, 0x5b, 0x30, 0x6d,
		0x1b, 0x5b, 0x33, 0x38, 0x3b, 0x35, 0x3b, 0x31, 0x32, 0x39, 0x6d, 0xe2, 0x96, 0x88, 0x1b, 0x5b, 0x30, 0x6d,
		0x1b, 0x5b, 0x33, 0x38, 0x3b, 0x35, 0x3b, 0x31, 0x39, 0x37, 0x6d, 0xe2, 0x96, 0x88, 0x1b, 0x5b, 0x30, 0x6d,
		0x1b, 0x5b, 0x33, 0x38, 0x3b, 0x35, 0x3b, 0x31, 0x36, 0x6d, 0xe2, 0x96, 0x88, 0x1b, 0x5b, 0x30, 0x6d,
		0x1b, 0x5b, 0x33, 0x38, 0x3b, 0x35, 0x3b, 0x31, 0x32, 0x39, 0x6d, 0xe2, 0x96, 0x88, 0x1b, 0x5b, 0x30, 0x6d,
	}, []byte(sb.String()))
}

func TestHeatNoColor(t *testing.T) {
	var sb strings.Builder

	HeatWriteLinear(&sb, 0, 0, 0)
	HeatWriteLinear(&sb, -5, 0, 10)
	HeatWriteLinear(&sb, 15, 0, 10)
	HeatWriteLinear(&sb, 0, 0, 10)
	HeatWriteLinear(&sb, 5, 0, 10)
	HeatWriteLinear(&sb, 10, 0, 10)
	HeatWriteLinear(&sb, 5, 5, 2)
	HeatWriteLinear(&sb, 750, 500, 1000)

	assert.Equal(t, "--9-59-5", sb.String())
}
