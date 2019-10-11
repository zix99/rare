package expressions

import (
	"fmt"
	"math"
	"strconv"
)

// KeyBuilderFunction defines a helper function at runtime
type KeyBuilderFunction func([]KeyBuilderStage) KeyBuilderStage

func kfBucket(args []KeyBuilderStage) KeyBuilderStage {
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		if len(args) != 2 {
			return ErrorArgCount
		}
		val, err := strconv.Atoi(args[0](context))
		if err != nil {
			return ErrorBucket
		}

		bucketSize, err := strconv.Atoi(args[1](context))
		if err != nil {
			return ErrorBucketSize
		}

		return strconv.Itoa((val / bucketSize) * bucketSize)
	})
}

var byteSizes = [...]string{"B", "KB", "MB", "GB", "TB", "PB"}

func kfBytesize(args []KeyBuilderStage) KeyBuilderStage {
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		if len(args) < 1 {
			return ErrorArgCount
		}
		val, err := strconv.Atoi(args[0](context))
		if err != nil {
			return ErrorType
		}

		labelIdx := 0
		for val >= 1024 && labelIdx < len(byteSizes)-1 {
			val = val / 1024
			labelIdx++
		}

		return fmt.Sprintf("%d %s", val, byteSizes[labelIdx])
	})
}

func kfExpBucket(args []KeyBuilderStage) KeyBuilderStage {
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		if len(args) != 1 {
			return ErrorArgCount
		}
		val, err := strconv.Atoi(args[0](context))
		if err != nil {
			return ErrorType
		}
		logVal := int(math.Log10(float64(val)))

		return strconv.Itoa(int(math.Pow10(logVal)))
	})
}

func arithmaticHelperi(equation func(int, int) int) KeyBuilderFunction {
	return KeyBuilderFunction(func(args []KeyBuilderStage) KeyBuilderStage {
		return KeyBuilderStage(func(context KeyBuilderContext) string {
			if len(args) < 2 {
				return ErrorArgCount
			}

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

var defaultFunctions = map[string]KeyBuilderFunction{
	"bucket":    KeyBuilderFunction(kfBucket),
	"expbucket": KeyBuilderFunction(kfExpBucket),
	"bytesize":  KeyBuilderFunction(kfBytesize),
	"sumi":      arithmaticHelperi(func(a, b int) int { return a + b }),
	"subi":      arithmaticHelperi(func(a, b int) int { return a - b }),
	"multi":     arithmaticHelperi(func(a, b int) int { return a * b }),
	"divi":      arithmaticHelperi(func(a, b int) int { return a / b }),
}
