package exprofiler

import (
	"testing"

	"github.com/zix99/rare/pkg/expressions"
	"github.com/zix99/rare/pkg/expressions/stdlib"

	"github.com/stretchr/testify/assert"
)

func TestExpressionStats(t *testing.T) {
	kb := stdlib.NewStdKeyBuilder()
	ckb, _ := kb.Compile("this is {0} a {test}")
	ctx := &expressions.KeyBuilderContextArray{}
	stats := GetMetrics(ckb, ctx)

	assert.Equal(t, 1, stats.MatchLookups)
	assert.Equal(t, 1, stats.MatchLookups)
}
