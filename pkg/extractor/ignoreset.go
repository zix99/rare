package extractor

import (
	"rare/pkg/expressions"
	"strings"
)

type IgnoreSet interface {
	IgnoreMatch(context expressions.KeyBuilderContext) bool
}

type ExpressionIgnoreSet struct {
	expressions []*expressions.CompiledKeyBuilder
}

func NewIgnoreExpressions(expSet ...string) (IgnoreSet, error) {
	if expSet == nil {
		return nil, nil
	}
	igSet := &ExpressionIgnoreSet{
		expressions: make([]*expressions.CompiledKeyBuilder, 0),
	}

	for _, exp := range expSet {
		compiled, err := expressions.NewKeyBuilder().Compile(exp)
		if err != nil {
			return nil, err
		}
		igSet.expressions = append(igSet.expressions, compiled)
	}

	return igSet, nil
}

func (s *ExpressionIgnoreSet) IgnoreMatch(context expressions.KeyBuilderContext) bool {
	if len(s.expressions) == 0 {
		return false
	}
	for _, exp := range s.expressions {
		result := strings.TrimSpace(exp.BuildKey(context))
		if expressions.Truthy(result) {
			return true
		}
	}

	return false
}
