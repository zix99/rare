package humanize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestH(t *testing.T) {
	assert.Equal(t, "Hello 1,000", H("Hello %d", 1000))
}

func TestHDisabled(t *testing.T) {
	Enabled = false
	assert.Equal(t, "Hello 1000", H("Hello %d", 1000))
	assert.Equal(t, "1000", Hi(1000))
	assert.Equal(t, "1000.000000", Hf(1000.0))
	assert.Equal(t, "1000.000000", Hfd(1000.0, 5))
	assert.Equal(t, "12341234", ByteSize(12341234))
	Enabled = true
}

func TestHi(t *testing.T) {
	assert.Equal(t, "1,500", Hi(1500))
}

func TestHf(t *testing.T) {
	assert.Equal(t, "1,234,567.8912", Hf(1234567.89121111))
}

func TestHfd(t *testing.T) {
	assert.Equal(t, "1,234,567.89", Hfd(1234567.89121111, 2))
}

func TestByteSize(t *testing.T) {
	assert.Equal(t, "123 B", ByteSize(123))
	assert.Equal(t, "1000 B", ByteSize(1000))
	assert.Equal(t, "1.46 KB", ByteSize(1500))
	assert.Equal(t, "2.00 MB", ByteSize(2*1024*1024))
	assert.Equal(t, "5.10 GB", ByteSize(5*1024*1024*1024+100*1024*1024))
	assert.Equal(t, "5 GB", AlwaysByteSize(5*1024*1024*1024+100*1024*1024, 0))
}

// 459.8 ns/op	      40 B/op	       3 allocs/op
func BenchmarkByteSize(b *testing.B) {
	Enabled = false
	for i := 0; i < b.N; i++ {
		AlwaysByteSize(5*1024*1024*1024+100*1024*1024, 2)
	}
}
