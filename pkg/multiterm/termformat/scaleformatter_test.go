package termformat

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPassthru(t *testing.T) {
	assert.Equal(t, "123", Passthru(123, 0, 10))
}

func TestDefault(t *testing.T) {
	assert.Equal(t, "1,234", Default(1234, 0, 10))
}

func TestFormatMapper(t *testing.T) {
	ts := ToFormatter(func(val int64) string {
		return strconv.FormatInt(val, 10)
	})

	assert.Equal(t, "123", ts(123, 0, 10))
}
