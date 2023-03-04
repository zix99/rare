package markdowncli

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// It's not a complete test, but it's better to at least do some code-path exercise
func TestMarkdown(t *testing.T) {
	r := strings.NewReader("# Hello World\nThis is a `string`\n\n```\nandsomecode\n```")
	w := &bytes.Buffer{}
	WriteMarkdownToBuf(w, r)
	assert.Equal(t, "# Hello World\n This is a `string`\n \n  andsomecode\n", w.String())
	fmt.Println(w.String())
}

func TestNoteBlock(t *testing.T) {
	r := strings.NewReader("# Title\n!!! note\n    this is a note block\n\n")
	w := &bytes.Buffer{}
	WriteMarkdownToBuf(w, r)

	assert.Equal(t, "# Title\n !!! note\n     this is a note block\n\n", w.String())
}
