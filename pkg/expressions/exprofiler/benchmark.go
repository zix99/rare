package exprofiler

import (
	"rare/pkg/expressions"
	"time"
)

func Benchmark(kb *expressions.CompiledKeyBuilder, ctx expressions.KeyBuilderContext) (duration time.Duration, iterations int) {
	const minTime = 500 * time.Millisecond
	iterations = 100_000

	for {
		start := time.Now()
		for i := 0; i < iterations; i++ {
			kb.BuildKey(ctx)
		}
		stop := time.Now()

		duration = stop.Sub(start)
		if duration >= minTime {
			return duration, iterations
		}

		iterations *= 4
	}
}
