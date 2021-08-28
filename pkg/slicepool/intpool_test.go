package slicepool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntPoolSimple(t *testing.T) {
	ip := NewIntPool(4)
	assert.Len(t, ip.Get(3), 3)
	assert.Len(t, ip.Get(1), 1)
	assert.Len(t, ip.Get(2), 2)
}

func TestIntPoolPanic(t *testing.T) {
	ip := NewIntPool(2)
	assert.Panics(t, func() {
		ip.Get(10)
	})
}

func TestIntPoolExact(t *testing.T) {
	ip := NewIntPool(2)
	assert.Len(t, ip.Get(2), 2)
}

var benchmarkSize = 10

// BenchmarkIntPool-4   	21709514	        49.45 ns/op	      80 B/op	       0 allocs/op
func BenchmarkIntPool(b *testing.B) {
	ip := NewIntPool(16 * 1024)
	for i := 0; i < b.N; i++ {
		v := ip.Get(benchmarkSize)
		fill(v)
		if len(v) != benchmarkSize {
			panic("Boom")
		}
	}
}

// BenchmarkRuntimeAlloc-4   	13441696	       106.5 ns/op	      80 B/op	       1 allocs/op
func BenchmarkRuntimeAlloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := make([]int, benchmarkSize)
		fill(v)
		if len(v) != benchmarkSize {
			panic("Boom")
		}
	}
}

func fill(b []int) {
	for i := 0; i < len(b); i++ {
		b[i] = i
	}
}
