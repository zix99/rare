package exprofiler

import (
	"rare/pkg/expressions"
	"rare/pkg/expressions/stdlib"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBenchmarking(t *testing.T) {
	kb := stdlib.NewStdKeyBuilder()
	ckb, _ := kb.Compile("hello {0}")
	ctx := expressions.KeyBuilderContextArray{}

	duration, iterations := Benchmark(ckb, &ctx)

	assert.NotZero(t, duration.Milliseconds())
	assert.NotZero(t, iterations)
}
