package grok

import (
	"rare/pkg/fastregex"
	"rare/pkg/grok/stdpat"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllStdlibPatternsCompile(t *testing.T) {
	grok := New()
	for _, k := range stdpat.Stdlib().Patterns() {
		pattern, err := grok.RewriteGrokPattern("%{" + k + "}")
		assert.NoError(t, err, k)
		_, err = fastregex.Compile(pattern)
		assert.NoError(t, err, k, pattern)
	}
}

func TestStdDoubleNested(t *testing.T) {
	grok := New()
	r, err := grok.RewriteGrokPattern("%{MAC:iph}")
	assert.NoError(t, err)
	assert.Equal(t, "(?P<iph>(?:(?:(?:[A-Fa-f0-9]{4}\\.){2}[A-Fa-f0-9]{4})|(?:(?:[A-Fa-f0-9]{2}-){5}[A-Fa-f0-9]{2})|(?:(?:[A-Fa-f0-9]{2}:){5}[A-Fa-f0-9]{2})))", r)
	_, err = fastregex.Compile(r)
	assert.NoError(t, err)
}

func TestStdNestedPattern(t *testing.T) {
	grok := New()
	r, err := grok.RewriteGrokPattern("User: %{USER:u}")
	assert.NoError(t, err)
	assert.Equal(t, "User: (?P<u>[a-zA-Z0-9._-]+)", r)
}
