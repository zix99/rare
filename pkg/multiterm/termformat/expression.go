package termformat

import (
	"rare/pkg/expressions"
	"rare/pkg/expressions/funclib"
	"rare/pkg/slicepool"
	"strconv"
)

type formatExpressionContext struct {
	val, min, max int64
}

var _ expressions.KeyBuilderContext = &formatExpressionContext{}

func (s *formatExpressionContext) GetMatch(idx int) string {
	switch idx {
	case 0:
		return strconv.FormatInt(s.val, 10)
	case 1:
		return strconv.FormatInt(s.min, 10)
	case 2:
		return strconv.FormatInt(s.max, 10)
	}
	return ""
}

func (s *formatExpressionContext) GetKey(key string) string {
	switch key {
	case "val", "value":
		return strconv.FormatInt(s.val, 10)
	case "min":
		return strconv.FormatInt(s.min, 10)
	case "max":
		return strconv.FormatInt(s.max, 10)
	}
	return ""
}

// Build a formatter using the default expression engine
func FromExpression(expr string) (Formatter, error) {
	kb, err := funclib.NewKeyBuilder().Compile(expr)
	if err != nil {
		return nil, err
	}

	pool := slicepool.NewObjectPool[formatExpressionContext](10)

	return func(val, min, max int64) string {
		ctx := pool.Get()
		defer pool.Return(ctx)

		*ctx = formatExpressionContext{val, min, max}
		return kb.BuildKey(ctx)
	}, nil
}

func MustFromExpression(expr string) Formatter {
	kb, err := FromExpression(expr)
	if err != nil {
		panic(err)
	}
	return kb
}
