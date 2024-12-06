//go:build linux && cgo && pcre2

package fastregex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaptureGroupCount(t *testing.T) {
	re, err := Compile("this is (.+) test (.+)")
	assert.NoError(t, err)
	assert.Equal(t, 3, re.CreateInstance().(*pcre2Regexp).GroupCount())
}
