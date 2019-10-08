package extractor

import (
	"strconv"
)

// KeyBuilderFunction defines a helper function at runtime
type KeyBuilderFunction func([]KeyBuilderStage) KeyBuilderStage

func kfBucket(args []KeyBuilderStage) KeyBuilderStage {
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val, err := strconv.Atoi(args[0](context))
		if err != nil {
			return "<BUCKET-ERROR>"
		}

		bucketSize, err := strconv.Atoi(args[1](context))
		if err != nil {
			return "<BUCKET-SIZE>"
		}

		return strconv.Itoa((val / bucketSize) * bucketSize)
	})
}

var defaultFunctions = map[string]KeyBuilderFunction{
	"bucket": KeyBuilderFunction(kfBucket),
}
