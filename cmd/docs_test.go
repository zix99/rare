package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocs(t *testing.T) {
	assert.NoError(t, testCommand(docsCommand(), ``))
	assert.NoError(t, testCommand(docsCommand(), `expressions`))
	assert.Error(t, testCommand(docsCommand(), `no-exist`))
}
