package aggregation

import (
	"rare/pkg/expressions"
	"rare/pkg/expressions/stdlib"
	"rare/pkg/stringSplitter"
)

// Context to build run expressions that use current value
type exprAccumulatorContext struct {
	current string
	match   string
	sub     expressions.KeyBuilderContext
}

func (s *exprAccumulatorContext) GetMatch(idx int) (ret string) {
	if idx == 0 {
		return s.match
	}

	// Index 1+, parse the string as if it's a range
	splitter := stringSplitter.Splitter{S: s.match, Delim: expressions.ArraySeparatorString}
	for i := 0; i < idx; i++ {
		ret = splitter.Next()
	}
	return
}

func (s *exprAccumulatorContext) GetKey(key string) string {
	if key == "." {
		return s.current
	}
	if s.sub != nil {
		return s.sub.GetKey(key)
	}
	return ""
}

// Basic accumulator that will sample a new value and accumulate it to a current `value`
type ExprAccumulator struct {
	expr  *expressions.CompiledKeyBuilder
	value string
}

var _ Aggregator = &ExprAccumulator{}

func NewExprAccumulator(expr, initial string) (*ExprAccumulator, error) {
	compiler := stdlib.NewStdKeyBuilder()
	cExpr, err := compiler.Compile(expr)
	if err != nil {
		return nil, err
	}

	return &ExprAccumulator{
		expr:  cExpr,
		value: initial,
	}, nil
}

func (s *ExprAccumulator) SampleEx(element string, subcontext expressions.KeyBuilderContext) {
	context := exprAccumulatorContext{
		current: s.value,
		match:   element,
		sub:     subcontext,
	}
	ret := s.expr.BuildKey(&context)
	if expressions.Truthy(ret) {
		s.value = ret
	}
}

func (s *ExprAccumulator) Sample(element string) {
	s.SampleEx(element, nil)
}

func (s *ExprAccumulator) ParseErrors() uint64 {
	return 0
}

func (s *ExprAccumulator) Value() string {
	return s.value
}

// Set of accumulators that will accumulate more than just a single value, and allow reference to each other

type ExprAccumulatorPair struct {
	Name  string
	Accum *ExprAccumulator
}

type ExprAccumulatorSet struct {
	accums []ExprAccumulatorPair
}

var _ Aggregator = &ExprAccumulatorSet{}

func NewExprAccumulatorSet() *ExprAccumulatorSet {
	return &ExprAccumulatorSet{}
}

func (s *ExprAccumulatorSet) Add(name, expr, initial string) error {
	accum, err := NewExprAccumulator(expr, initial)
	if err != nil {
		return err
	}

	s.accums = append(s.accums, ExprAccumulatorPair{
		name,
		accum,
	})

	return nil
}

func (s *ExprAccumulatorSet) Sample(element string) {
	for _, accum := range s.accums {
		accum.Accum.SampleEx(element, s)
	}
}

func (s *ExprAccumulatorSet) ParseErrors() uint64 {
	return 0
}

func (s *ExprAccumulatorSet) Items() []ExprAccumulatorPair {
	return s.accums
}

func (s *ExprAccumulatorSet) GetKey(key string) string {
	for _, ele := range s.accums {
		if ele.Name == key {
			return ele.Accum.value
		}
	}
	return ""
}

func (s *ExprAccumulatorSet) GetMatch(idx int) string {
	return ""
}
