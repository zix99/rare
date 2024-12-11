package cmd

import (
	"bytes"
	"rare/docs"
	"rare/pkg/markdowncli"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocs(t *testing.T) {
	assert.NoError(t, testCommand(docsCommand(), ``))
	assert.NoError(t, testCommand(docsCommand(), `-n expressions`))
	assert.NoError(t, testCommand(docsCommand(), `-n exp`)) // partial name
	assert.Error(t, testCommand(docsCommand(), `no-exist`))
}

// Test that all the docs make it through markdown parser with a result
func TestDocsProcess(t *testing.T) {
	entries, _ := docs.DocFS.ReadDir(docs.BasePath)
	assert.NotZero(t, len(entries))

	for _, entry := range entries {
		entry := entry
		t.Run(entry.Name(), func(t *testing.T) {
			t.Parallel()
			f, err := docs.DocFS.Open(docs.BasePath + "/" + entry.Name())
			assert.NoError(t, err)
			defer f.Close()

			var buf bytes.Buffer
			markdowncli.WriteMarkdownToBuf(&buf, f)

			assert.NotZero(t, buf.Len())
		})
	}
}
