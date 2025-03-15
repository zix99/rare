package stdlib

import (
	"rare/pkg/expressions"
	"strconv"
)

type typedStage[T any] func(context expressions.KeyBuilderContext) (val T, ok bool)

type typedStageParser[T any] func(string) (T, bool)

var (
	typedParserFloat = func(s string) (float64, bool) {
		if v, err := strconv.ParseFloat(s, 64); err == nil {
			return v, true
		}
		return 0.0, false
	}
	typedParsedInt = func(s string) (int, bool) {
		if v, err := strconv.Atoi(s); err == nil {
			return v, true
		}
		return 0, false
	}
)

func typedLiteral[T any](val T) typedStage[T] {
	return func(context expressions.KeyBuilderContext) (T, bool) {
		return val, true
	}
}

func evalDynamicStage[T any](stage expressions.KeyBuilderStage, parser typedStageParser[T]) (typedStage[T], bool) {
	if val, ok := expressions.EvalStaticStage(stage); ok {
		if pval, ok := parser(val); ok {
			return func(context expressions.KeyBuilderContext) (val T, ok bool) {
				return pval, true
			}, true
		} else {
			return nil, false
		}
	}

	return func(context expressions.KeyBuilderContext) (T, bool) {
		return parser(stage(context))
	}, true
}

func mapDynamicArgs[T any](args []expressions.KeyBuilderStage, parser typedStageParser[T]) ([]typedStage[T], bool) {
	ret := make([]typedStage[T], len(args))

	for i, arg := range args {
		var ok bool
		ret[i], ok = evalDynamicStage(arg, parser)
		if !ok {
			return nil, false
		}
	}

	return ret, true
}
