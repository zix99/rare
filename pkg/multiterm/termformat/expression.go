package termformat

import (
	"rare/pkg/expressions"
	"rare/pkg/expressions/funclib"
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

// Special case of the expression builder where if a function exists
// it will be used to format. eg. providing just `bytesize` will yield `{bytesize {0}}`
// This works well, since you'd probably never intend to return the word "bytesize" for the format
func expandCompileExpression(expr string) (*expressions.CompiledKeyBuilder, *expressions.CompilerErrors) {
	if funclib.FunctionExists(expr) {
		expr = "{" + expr + " {0}}"
	}
	return funclib.NewKeyBuilder().Compile(expr)
}

// Build a formatter using the default expression engine
// Single-threaded use only
func FromExpression(expr string) (Formatter, error) {
	kb, err := expandCompileExpression(expr)
	if err != nil {
		return nil, err
	}

	ctx := &formatExpressionContext{}
	return func(val, min, max int64) string {
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
