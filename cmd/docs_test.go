package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocs(t *testing.T) {
	assert.NoError(t, testCommand(docsCommand(), ``))
	assert.NoError(t, testCommand(docsCommand(), `-n expressions`))
	assert.NoError(t, testCommand(docsCommand(), `-n exp`))
	assert.Error(t, testCommand(docsCommand(), `no-exist`))
}
