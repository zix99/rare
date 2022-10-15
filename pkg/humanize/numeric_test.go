package humanize

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatInt(t *testing.T) {
	assert.Equal(t, "1", humanizeInt(1))
	assert.Equal(t, "10", humanizeInt(10))
	assert.Equal(t, "100", humanizeInt(100))
	assert.Equal(t, "1,000", humanizeInt(1000))
	assert.Equal(t, "10,000", humanizeInt(10000))

	assert.Equal(t, "-100", humanizeInt(-100))
	assert.Equal(t, "-1,000", humanizeInt(-1000))
	assert.Equal(t, "-123,123", humanizeInt(-123123))
}

func TestFormatFloat(t *testing.T) {
	assert.Equal(t, "1", humanizeFloat(1.0, 0))
	assert.Equal(t, "12", humanizeFloat(12.0, 0))
	assert.Equal(t, "123", humanizeFloat(123.0, 0))
	assert.Equal(t, "1,234", humanizeFloat(1234.0, 0))
	assert.Equal(t, "12,345.0", humanizeFloat(12345.0, 1))
	assert.Equal(t, "112,345.0", humanizeFloat(112345.0, 1))
	assert.Equal(t, "1", humanizeFloat(1.123, 0))
	assert.Equal(t, "-1", humanizeFloat(-1.123, 0))
	assert.Equal(t, "1,123,123", humanizeFloat(1123123.123, 0))
	assert.Equal(t, "-1,123,123", humanizeFloat(-1123123.123, 0))
	assert.Equal(t, "1,123,123.12", humanizeFloat(1123123.123, 2))
	assert.Equal(t, "1,123,123.123456", humanizeFloat(1123123.123456, 6))
	assert.Equal(t, "-1,123,123.123456", humanizeFloat(-1123123.123456, 6))
	assert.Equal(t, "-111,121,231,233,123.125000", humanizeFloat(-111121231233123.123456, 6))
	assert.Equal(t, "111,121,231,233,123.125000", humanizeFloat(111121231233123.123456, 6))
}

func BenchmarkFormatInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		humanizeInt(10000)
	}
}

func BenchmarkItoa(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.Itoa(10000)
	}
}

func BenchmarkFormatFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		humanizeFloat(10000.0, 2)
	}
}

func BenchmarkFloatSprintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("%02f", 10000.0)
	}
}
