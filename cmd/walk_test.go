package cmd

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyWalkCommand(t *testing.T) {
	o, e, err := testCommandCapture(walkCommand(), "")
	assert.Error(t, err)
	assert.Equal(t, "", o)
	assert.Equal(t, "No paths found", e)
}

func TestWalkTestDataGlob(t *testing.T) {
	o, e, err := testCommandCapture(walkCommand(), "testdata/*.gz")
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join("testdata", "log.txt.gz")+"\n", o)
	assert.Equal(t, "Found 1 path(s)\n", e)
}

func TestWalkTestDataRecursive(t *testing.T) {
	o, e, err := testCommandCapture(walkCommand(), "-R testdata/")
	assert.NoError(t, err)
	assert.Contains(t, o, "log.txt")
	assert.Equal(t, "Found 7 path(s)\n", e)
}
