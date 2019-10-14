package extractor

import (
	"rare/pkg/expressions"
	"strings"
)

type IgnoreSet interface {
	IgnoreMatch(matchSet ...string) bool
}

type ExpressionIgnoreSet struct {
	expressions []*expressions.CompiledKeyBuilder
}

func NewIgnoreExpressions(expSet []string) (IgnoreSet, error) {
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

func (s *ExpressionIgnoreSet) IgnoreMatch(matchSet ...string) bool {
	if len(matchSet) == 0 || len(s.expressions) == 0 {
		return false
	}
	context := expressions.KeyBuilderContextArray{
		Elements: matchSet,
	}
	for _, exp := range s.expressions {
		result := strings.TrimSpace(exp.BuildKey(&context))
		if expressions.Truthy(result) {
			return true
		}
	}

	return false
}
