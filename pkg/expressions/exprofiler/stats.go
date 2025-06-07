package exprofiler

import "github.com/zix99/rare/pkg/expressions"

type ExpressionMetrics struct {
	MatchLookups, KeyLookups int
}

type trackingExpressionContext struct {
	Nested       expressions.KeyBuilderContext
	MatchLookups int
	KeyLookups   int
}

var _ expressions.KeyBuilderContext = &trackingExpressionContext{}

func (s *trackingExpressionContext) GetMatch(idx int) string {
	s.MatchLookups++
	return s.Nested.GetMatch(idx)
}

func (s *trackingExpressionContext) GetKey(key string) string {
	s.KeyLookups++
	return s.Nested.GetKey(key)
}

func GetMetrics(kb *expressions.CompiledKeyBuilder, ctx expressions.KeyBuilderContext) ExpressionMetrics {
	trackingContext := trackingExpressionContext{ctx, 0, 0}
	kb.BuildKey(&trackingContext)
	return ExpressionMetrics{
		MatchLookups: trackingContext.MatchLookups,
		KeyLookups:   trackingContext.KeyLookups,
	}
}
