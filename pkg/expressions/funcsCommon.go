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
		return stageError(ErrorArgCount)
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

func kfExpBucket(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 1 {
		return stageError(ErrorArgCount)
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
