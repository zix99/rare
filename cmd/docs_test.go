package cmd

import (
	"bufio"
	"bytes"
	"io"
	"rare/docs"
	"rare/pkg/markdowncli"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const MAX_LINE_LEN = 140

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
		t.Run("md:"+entry.Name(), func(t *testing.T) {
			t.Parallel()
			f, err := docs.DocFS.Open(docs.BasePath + "/" + entry.Name())
			assert.NoError(t, err)
			defer f.Close()

			var buf bytes.Buffer
			markdowncli.WriteMarkdownToBuf(&buf, f)

			assert.NotZero(t, buf.Len())
		})
		t.Run("len:"+entry.Name(), func(t *testing.T) {
			t.Parallel()
			f, err := docs.DocFS.Open(docs.BasePath + "/" + entry.Name())
			assert.NoError(t, err)
			defer f.Close()

			validateTerminalFit(t, f)
		})
	}
}

func validateTerminalFit(t *testing.T, r io.Reader) {
	scanner := bufio.NewScanner(r)

	ln := 0
	codeblock := false

	for scanner.Scan() {
		line := scanner.Text()
		ln++

		if strings.Contains(line, "```") {
			codeblock = !codeblock
		} else if !codeblock && len(line) > MAX_LINE_LEN {
			t.Errorf("Line %d too long (%d): %s", ln, len(line), line)
		}
	}
}
