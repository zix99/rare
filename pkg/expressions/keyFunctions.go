package expressions

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// KeyBuilderFunction defines a helper function at runtime
type KeyBuilderFunction func([]KeyBuilderStage) KeyBuilderStage

func kfCoalesce(args []KeyBuilderStage) KeyBuilderStage {
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		for _, arg := range args {
			val := arg(context)
			if val != "" {
				return val
			}
		}
		return ""
	})
}

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

func Truthy(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	return true
}

func kfNot(args []KeyBuilderStage) KeyBuilderStage {
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		if len(args) != 1 {
			return ErrorArgCount
		}
		if Truthy(args[0](context)) {
			return ""
		}
		return "1"
	})
}

// Simple helper that will take 2 or more integers, and apply an operation
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

// Checks equality, and returns truthy if equals, and empty if not
func arithmaticEqualityHelper(test func(int, int) bool) KeyBuilderFunction {
	return KeyBuilderFunction(func(args []KeyBuilderStage) KeyBuilderStage {
		return KeyBuilderStage(func(context KeyBuilderContext) string {
			if len(args) != 2 {
				return ErrorArgCount
			}

			left, err := strconv.Atoi(args[0](context))
			if err != nil {
				return ErrorType
			}
			right, err := strconv.Atoi(args[1](context))
			if err != nil {
				return ErrorType
			}

			if test(left, right) {
				return "1"
			}
			return ""
		})
	})
}

func stringHelper(equation func(string, string) string) KeyBuilderFunction {
	return KeyBuilderFunction(func(args []KeyBuilderStage) KeyBuilderStage {
		return KeyBuilderStage(func(context KeyBuilderContext) string {
			if len(args) < 2 {
				return ErrorArgCount
			}

			val := args[0](context)
			for i := 1; i < len(args); i++ {
				val = equation(val, args[i](context))
			}

			return val
		})
	})
}

var defaultFunctions = map[string]KeyBuilderFunction{
	"coalesce":  KeyBuilderFunction(kfCoalesce),
	"bucket":    KeyBuilderFunction(kfBucket),
	"expbucket": KeyBuilderFunction(kfExpBucket),
	"bytesize":  KeyBuilderFunction(kfBytesize),
	"sumi":      arithmaticHelperi(func(a, b int) int { return a + b }),
	"subi":      arithmaticHelperi(func(a, b int) int { return a - b }),
	"multi":     arithmaticHelperi(func(a, b int) int { return a * b }),
	"divi":      arithmaticHelperi(func(a, b int) int { return a / b }),
	"eq": stringHelper(func(a, b string) string {
		if a == b {
			return a
		}
		return ""
	}),
	"neq": stringHelper(func(a, b string) string {
		if a != b {
			return a
		}
		return ""
	}),
	"not": KeyBuilderFunction(kfNot),
	"lt":  arithmaticEqualityHelper(func(a, b int) bool { return a < b }),
	"gt":  arithmaticEqualityHelper(func(a, b int) bool { return a > b }),
	"lte": arithmaticEqualityHelper(func(a, b int) bool { return a <= b }),
	"gte": arithmaticEqualityHelper(func(a, b int) bool { return a >= b }),
}
