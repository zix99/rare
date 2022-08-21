//go:build !urfave_cli_no_docs

package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenDocs(t *testing.T) {
	out, _, err := testCommandCapture(gendocCommand(), ``)
	assert.NoError(t, err)
	assert.NotEmpty(t, out)

	out, _, err = testCommandCapture(gendocCommand(), `--man`)
	assert.NoError(t, err)
	assert.NotEmpty(t, out)
}
