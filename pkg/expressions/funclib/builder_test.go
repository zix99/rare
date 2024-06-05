package funclib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKeyBuilder(t *testing.T) {
	assert.NotNil(t, NewKeyBuilder())
}
