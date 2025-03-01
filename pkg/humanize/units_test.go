package humanize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByteSize(t *testing.T) {
	assert.Equal(t, "0 B", ByteSize(0))
	assert.Equal(t, "123 B", ByteSize(123))
	assert.Equal(t, "1000 B", ByteSize(1000))
	assert.Equal(t, "1.46 KB", ByteSize(1500))
	assert.Equal(t, "2.00 MB", ByteSize(2*1024*1024))
	assert.Equal(t, "5.10 GB", ByteSize(5*1024*1024*1024+100*1024*1024))
	assert.Equal(t, "5 GB", AlwaysByteSize(5*1024*1024*1024+100*1024*1024, 0))
	assert.Equal(t, "7 EB", AlwaysByteSize(7*1024*1024*1024*1024*1024*1024, 0))
}

func TestByteSizeSi(t *testing.T) {
	assert.Equal(t, "0 b", ByteSizeSi(0))
	assert.Equal(t, "123 b", ByteSizeSi(123))
	assert.Equal(t, "1.00 kB", ByteSizeSi(1000))
	assert.Equal(t, "1.50 kB", ByteSizeSi(1500))
	assert.Equal(t, "2.10 mB", ByteSizeSi(2*1024*1024))
	assert.Equal(t, "5.47 gB", ByteSizeSi(5*1024*1024*1024+100*1024*1024))
	assert.Equal(t, "5.5 gB", AlwaysByteSizeSi(5*1024*1024*1024+100*1024*1024, 1))
	assert.Equal(t, "7 pB", AlwaysByteSizeSi(7*1000*1000*1000*1000*1000, 0))
}

func TestDownscale(t *testing.T) {
	assert.Equal(t, "0", Downscale(0, 0))
	assert.Equal(t, "900", AlwaysDownscale(900, 0))
	assert.Equal(t, "10k", AlwaysDownscale(10000, 0))
	assert.Equal(t, "12k", AlwaysDownscale(12345, 0))
	assert.Equal(t, "12.35k", AlwaysDownscale(12345, 2))
	assert.Equal(t, "52.12M", Downscale(52_123_123, 2))
	assert.Equal(t, "3000T", AlwaysDownscale(3_000_000_000_000_000, 0))

	assert.Equal(t, "-900", AlwaysDownscale(-900, 0))
	assert.Equal(t, "-10k", AlwaysDownscale(-10000, 0))
	assert.Equal(t, "-12k", AlwaysDownscale(-12345, 0))
	assert.Equal(t, "-12.35k", AlwaysDownscale(-12345, 2))
	assert.Equal(t, "-52.12M", AlwaysDownscale(-52_123_123, 2))
}

// BenchmarkByteSize-14    	 5101810	       211.4 ns/op	       8 B/op	       1 allocs/op
func BenchmarkByteSize(b *testing.B) {
	Enabled = false
	for i := 0; i < b.N; i++ {
		AlwaysByteSize(5*1024*1024*1024+100*1024*1024, 2)
	}
}
