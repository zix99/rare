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
		return KeyBuilderStage(func(context KeyBuilderContext) string {
			final, err := strconv.Atoi(args[0](context))
			if err != nil {
				return ErrorNum
			}

			for i := 1; i < len(args); i++ {
				val, err := strconv.Atoi(args[i](context))
				if err != nil {
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
		return KeyBuilderStage(func(context KeyBuilderContext) string {
			final, err := strconv.ParseFloat(args[0](context), 64)
			if err != nil {
				return ErrorNum
			}

			for i := 1; i < len(args); i++ {
				val, err := strconv.ParseFloat(args[i](context), 64)
				if err != nil {
					return ErrorNum
				}
				final = equation(final, val)
			}

			return strconv.FormatFloat(final, 'f', -1, 64)
		}), nil
	})
}
