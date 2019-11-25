package expressions

import "strconv"

// Simple helper that will take 2 or more integers, and apply an operation
func arithmaticHelperi(equation func(int, int) int) KeyBuilderFunction {
	return KeyBuilderFunction(func(args []KeyBuilderStage) KeyBuilderStage {
		if len(args) < 2 {
			return stageLiteral(ErrorArgCount)
		}
		return KeyBuilderStage(func(context KeyBuilderContext) string {
			final, err := strconv.Atoi(args[0](context))
			if err != nil {
				return ErrorType
			}

			for i := 1; i < len(args); i++ {
				val, err := strconv.Atoi(args[i](context))
				if err != nil {
					return ErrorType
				}
				final = equation(final, val)
			}

			return strconv.Itoa(final)
		})
	})
}

// Simple helper that will take 2 or more integers, and apply an operation
func arithmaticHelperf(equation func(float64, float64) float64) KeyBuilderFunction {
	return KeyBuilderFunction(func(args []KeyBuilderStage) KeyBuilderStage {
		if len(args) < 2 {
			return stageLiteral(ErrorArgCount)
		}
		return KeyBuilderStage(func(context KeyBuilderContext) string {
			final, err := strconv.ParseFloat(args[0](context), 64)
			if err != nil {
				return ErrorType
			}

			for i := 1; i < len(args); i++ {
				val, err := strconv.ParseFloat(args[i](context), 64)
				if err != nil {
					return ErrorType
				}
				final = equation(final, val)
			}

			return strconv.FormatFloat(final, 'f', -1, 64)
		})
	})
}
