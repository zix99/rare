package stdlib

import (
	. "rare/pkg/expressions" //lint:ignore ST1001 Legacy
	"strconv"
)

// Simple helper that will take 2 or more integers, and apply an operation
func arithmaticHelperi(equation func(int, int) int) KeyBuilderFunction {
	return KeyBuilderFunction(func(args []KeyBuilderStage) (KeyBuilderStage, error) {
		if len(args) < 2 {
			return stageErrArgRange(args, "2+")
		}

		typedArgs, tOk := mapTypedArgs(args, typedParserInt)
		if !tOk {
			return stageError(ErrNum)
		}

		return KeyBuilderStage(func(context KeyBuilderContext) string {
			final, ok := typedArgs[0](context)
			if !ok {
				return ErrorNum
			}

			for i := 1; i < len(args); i++ {
				val, ok := typedArgs[i](context)
				if !ok {
					return ErrorNum
				}
				final = equation(final, val)
			}

			return strconv.Itoa(final)
		}), nil
	})
}

// Simple helper that will take 2 or more integers, and apply an operation
func arithmaticHelperf(equation func(float64, float64) float64) KeyBuilderFunction {
	return KeyBuilderFunction(func(args []KeyBuilderStage) (KeyBuilderStage, error) {
		if len(args) < 2 {
			return stageErrArgRange(args, "2+")
		}

		typedArgs, tOk := mapTypedArgs(args, typedParserFloat)
		if !tOk {
			return stageError(ErrNum)
		}

		return KeyBuilderStage(func(context KeyBuilderContext) string {
			final, ok := typedArgs[0](context)
			if !ok {
				return ErrorNum
			}

			for i := 1; i < len(args); i++ {
				val, ok := typedArgs[i](context)
				if !ok {
					return ErrorNum
				}
				final = equation(final, val)
			}

			return strconv.FormatFloat(final, 'f', -1, 64)
		}), nil
	})
}

// Helper that takes in a float, operates on it
func unaryArithmaticHelperf(op func(float64) float64) KeyBuilderFunction {
	return func(args []KeyBuilderStage) (KeyBuilderStage, error) {
		if len(args) != 1 {
			return stageErrArgCount(args, 1)
		}

		return func(context KeyBuilderContext) string {
			val, err := strconv.ParseFloat(args[0](context), 64)
			if err != nil {
				return ErrorNum
			}

			return strconv.FormatFloat(op(val), 'f', -1, 64)
		}, nil
	}
}

// Helper that takes in a float, operates on it, and spits out an int
func unaryArithmaticHelperfi(op func(float64) int64) KeyBuilderFunction {
	return func(args []KeyBuilderStage) (KeyBuilderStage, error) {
		if len(args) != 1 {
			return stageErrArgCount(args, 1)
		}

		return func(context KeyBuilderContext) string {
			val, err := strconv.ParseFloat(args[0](context), 64)
			if err != nil {
				return ErrorNum
			}

			return strconv.FormatInt(op(val), 10)
		}, nil
	}
}

// {round <val> [precision=0]}
func kfRound(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if !isArgCountBetween(args, 1, 2) {
		return stageErrArgRange(args, "1-2")
	}

	precision, precisionOk := EvalArgInt(args, 1, 0)
	if !precisionOk {
		return stageArgError(ErrConst, 1)
	}

	return func(context KeyBuilderContext) string {
		val, err := strconv.ParseFloat(args[0](context), 64)
		if err != nil {
			return ErrorNum
		}
		return strconv.FormatFloat(val, 'f', precision, 64)
	}, nil
}
