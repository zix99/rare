package expressions

import (
	"math"
	"strconv"
)

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
	if len(args) != 2 {
		return stageLiteral(ErrorArgCount)
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
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

func kfClamp(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 3 {
		return stageLiteral(ErrorArgCount)
	}

	min, minErr := strconv.Atoi(args[1](nil))
	max, maxErr := strconv.Atoi(args[2](nil))

	if minErr != nil || maxErr != nil {
		return stageLiteral(ErrorType)
	}

	return KeyBuilderStage(func(context KeyBuilderContext) string {
		arg0 := args[0](context)
		val, err := strconv.Atoi(arg0)
		if err != nil {
			return ErrorType
		}

		if val < min {
			return "min"
		} else if val > max {
			return "max"
		} else {
			return arg0
		}
	})
}

func kfExpBucket(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 1 {
		return stageLiteral(ErrorArgCount)
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val, err := strconv.Atoi(args[0](context))
		if err != nil {
			return ErrorType
		}
		logVal := int(math.Log10(float64(val)))

		return strconv.Itoa(int(math.Pow10(logVal)))
	})
}
