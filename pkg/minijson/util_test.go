package minijson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapToJSON(t *testing.T) {
	m := map[string]string{
		"a": "b",
		"c": "123",
	}
	val := MarshalStringMapInferred(m)
	assert.Equal(t, `{"a": "b", "c": "123"}`, val)
}
