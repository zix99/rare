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
