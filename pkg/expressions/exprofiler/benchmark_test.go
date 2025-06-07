package exprofiler

import (
	"testing"

	"github.com/zix99/rare/pkg/expressions"
	"github.com/zix99/rare/pkg/expressions/stdlib"

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
