package markdowncli

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// It's not a complete test, but it's better to at least do some code-path exercise
func TestMarkdown(t *testing.T) {
	r := strings.NewReader("# Hello World\nThis is a `string`\n\n```\nandsomecode\n```")
	w := &bytes.Buffer{}
	WriteMarkdownToBuf(w, r)
	fmt.Println(w.String())
}
