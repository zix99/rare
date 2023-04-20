package stdlib

import (
	"math"
	"strconv"

	. "rare/pkg/expressions" //lint:ignore ST1001 Legacy
)

func kfCoalesce(args []KeyBuilderStage) (KeyBuilderStage, error) {
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		for _, arg := range args {
			val := arg(context)
			if val != "" {
				return val
			}
		}
		return ""
	}), nil
}

func kfBucket(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 2 {
		return stageError(ErrArgCount)
	}

	bucketSize, err := strconv.Atoi(EvalStageOrDefault(args[1], ""))
	if err != nil {
		return stageError(ErrNum)
	}

	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val, err := strconv.Atoi(args[0](context))
		if err != nil {
			return ErrorNum
		}

		return strconv.Itoa((val / bucketSize) * bucketSize)
	}), nil
}

func kfClamp(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 3 {
		return stageError(ErrArgCount)
	}

	min, minErr := strconv.Atoi(EvalStageOrDefault(args[1], ""))
	max, maxErr := strconv.Atoi(EvalStageOrDefault(args[2], ""))

	if minErr != nil || maxErr != nil {
		return stageError(ErrNum)
	}

	return KeyBuilderStage(func(context KeyBuilderContext) string {
		arg0 := args[0](context)
		val, err := strconv.Atoi(arg0)
		if err != nil {
			return ErrorNum
		}

		if val < min {
			return "min"
		} else if val > max {
			return "max"
		} else {
			return arg0
		}
	}), nil
}

func kfExpBucket(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 1 {
		return stageError(ErrArgCount)
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val, err := strconv.Atoi(args[0](context))
		if err != nil {
			return ErrorNum
		}
		logVal := int(math.Log10(float64(val)))

		return strconv.Itoa(int(math.Pow10(logVal)))
	}), nil
}
