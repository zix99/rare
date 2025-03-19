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
		return stageErrArgCount(args, 2)
	}

	bucketSize, bucketSizeOk := EvalStageInt64(args[1])
	if !bucketSizeOk {
		return stageArgError(ErrNum, 1)
	}
	if bucketSize <= 0 {
		return stageArgError(ErrValue, 1)
	}

	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val, err := strconv.ParseInt(args[0](context), 10, 64)
		if err != nil {
			return ErrorNum
		}

		bucket := (val / bucketSize) * bucketSize
		if val < 0 {
			bucket -= bucketSize
		}

		return strconv.FormatInt(bucket, 10)
	}), nil
}

func kfBucketRange(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 2 {
		return stageErrArgCount(args, 2)
	}

	bucketSize, bucketSizeOk := EvalStageInt64(args[1])
	if !bucketSizeOk {
		return stageArgError(ErrNum, 1)
	}
	if bucketSize <= 0 {
		return stageArgError(ErrValue, 1)
	}

	return func(context KeyBuilderContext) string {
		val, err := strconv.ParseInt(args[0](context), 10, 64)
		if err != nil {
			return ErrorNum
		}

		var start, end int64
		start = (val / bucketSize) * bucketSize
		if val < 0 {
			start -= bucketSize
		}
		end = start + (bucketSize - 1)

		ret := make([]byte, 0, 20)
		ret = strconv.AppendInt(ret, start, 10)
		ret = append(ret, " - "...)
		ret = strconv.AppendInt(ret, end, 10)
		return string(ret)
	}, nil
}

func kfClamp(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 3 {
		return stageErrArgCount(args, 3)
	}

	min, minOk := EvalStageInt(args[1])
	max, maxOk := EvalStageInt(args[2])

	if !minOk {
		return stageArgError(ErrNum, 1)
	}
	if !maxOk {
		return stageArgError(ErrNum, 2)
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
		return stageErrArgCount(args, 1)
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
