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
	typedParserInt = func(s string) (int, bool) {
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

// Using a parser, turn a stage into a typed-stage. Return false on error (unable to parse static)
// Returns a literal if able, otherwise returns a parser-wrapper
// Useful when an expression function uses more than one typed argument and likely won't be optimized-out on its own
// Saves the parse time (generally 10-50%)
func evalTypedStage[T any](stage expressions.KeyBuilderStage, parser typedStageParser[T]) (typedStage[T], bool) {
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

// Map an argument slice to typed stages
func mapTypedArgs[T any](args []expressions.KeyBuilderStage, parser typedStageParser[T]) ([]typedStage[T], bool) {
	ret := make([]typedStage[T], len(args))

	for i, arg := range args {
		var ok bool
		ret[i], ok = evalTypedStage(arg, parser)
		if !ok {
			return nil, false
		}
	}

	return ret, true
}
