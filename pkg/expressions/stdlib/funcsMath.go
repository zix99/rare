package stdlib

import (
	"strconv"
	"strings"

	"github.com/zix99/rare/pkg/expressions"
	"github.com/zix99/rare/pkg/expressions/stdmath"
	"github.com/zix99/rare/pkg/slicepool"
)

type keyBuilderContextWrapper struct {
	sub    expressions.KeyBuilderContext
	errors int
}

func (s *keyBuilderContextWrapper) GetMatch(idx int) float64 {
	val := s.sub.GetMatch(idx)
	if f, err := strconv.ParseFloat(val, 64); err == nil {
		return f
	}
	s.errors++
	return 0
}

func (s *keyBuilderContextWrapper) GetKey(key string) float64 {
	val := s.sub.GetKey(key)
	if f, err := strconv.ParseFloat(val, 64); err == nil {
		return f
	}
	s.errors++
	return 0
}

func kfMath(args []expressions.KeyBuilderStage) (expressions.KeyBuilderStage, error) {
	// Collapse all arguments to a single expression
	var sb strings.Builder
	for i, arg := range args {
		s, ok := expressions.EvalStaticStage(arg)
		if !ok {
			return stageArgError(ErrConst, i)
		}
		sb.WriteString(s)
	}

	// Compile
	expr, err := stdmath.Compile(sb.String())
	if err != nil {
		return stageErrorf(ErrParsing, err.Error())
	}

	// Runner
	ctxPool := slicepool.NewObjectPool[keyBuilderContextWrapper](5)

	return func(ctx expressions.KeyBuilderContext) string {
		mathCtx := ctxPool.Get()
		defer ctxPool.Return(mathCtx)

		*mathCtx = keyBuilderContextWrapper{
			sub:    ctx,
			errors: 0,
		}

		val := expr.Eval(mathCtx)

		if mathCtx.errors > 0 {
			return ErrorNum
		}

		return strconv.FormatFloat(val, 'f', -1, 64)
	}, nil
}
