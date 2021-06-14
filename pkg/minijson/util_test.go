package minijson

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapToJSON(t *testing.T) {
	m := map[string]string{
		"a": "b",
		"c": "123",
	}
	val := MarshalStringMapInferred(m)

	// Order is non-deterministic, so check for contents
	assert.Contains(t, val, `"a": "b"`)
	assert.Contains(t, val, `"c": "123"`)
	assert.True(t, strings.HasPrefix(val, "{"))
	assert.True(t, strings.HasSuffix(val, "}"))
}
