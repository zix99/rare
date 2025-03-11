package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildFormatter(t *testing.T) {
	f, err := BuildFormatter("{bytesize {0}}")
	assert.NotNil(t, f)
	assert.NoError(t, err)

	f, err = BuildFormatter("{bad expr")
	assert.Nil(t, f)
	assert.Error(t, err)

	f, err = BuildFormatter("") // default
	assert.NotNil(t, f)
	assert.NoError(t, err)
}

func TestBuildFormatterOrFail(t *testing.T) {
	assert.NotNil(t, BuildFormatterOrFail(""))
	assert.NotNil(t, BuildFormatterOrFail("{bytesize {0}}"))
	testLogFatal(t, 2, func() {
		BuildFormatterOrFail("{bad")
	})
}
