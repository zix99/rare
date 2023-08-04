package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildScaler(t *testing.T) {
	_, err := BuildScaler("")
	assert.NoError(t, err)

	_, err = BuildScaler("log10")
	assert.NoError(t, err)

	_, err = BuildScaler("bad-data")
	assert.Error(t, err)
}

func TestBuildScalerOrFail(t *testing.T) {
	assert.NotNil(t, BuildScalerOrFail(""))
	assert.NotNil(t, BuildScalerOrFail("linear"))
	testLogFatal(t, 2, func() {
		BuildScalerOrFail("fake")
	})
}
