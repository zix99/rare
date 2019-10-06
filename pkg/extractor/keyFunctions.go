package extractor

import (
	"fmt"
	"strconv"
)

func kfError(msg string) KeyBuilderStage {
	errMessage := fmt.Sprintf("<%s>", msg)
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		return errMessage
	})
}

func kfBucket(args []string) KeyBuilderStage {
	index, err1 := strconv.Atoi(args[0])
	bucketSize, err2 := strconv.Atoi(args[1])

	if err1 != nil || err2 != nil {
		return kfError("Invalid bucket arg")
	}

	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val, err := strconv.Atoi(context.GetMatch(index))
		if err != nil {
			return "<BUCKET-ERROR>"
		}
		return strconv.Itoa((val / bucketSize) * bucketSize)
	})
}

var defaultFunctions = map[string]KeyBuilderFunction{
	"bucket": KeyBuilderFunction(kfBucket),
}
