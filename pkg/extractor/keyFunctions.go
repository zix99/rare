package extractor

import (
	"strconv"
)

func kfBucket(args []string) KeyBuilderStage {
	index, err1 := strconv.Atoi(args[0])
	bucketSize, err2 := strconv.Atoi(args[1])

	if err1 != nil || err2 != nil {
		return stageError("Invalid bucket arg")
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
