package extractor

import "strconv"

func kfBucket(args []string) KeyBuilderStage {
	index, _ := strconv.Atoi(args[0])
	bucketSize, _ := strconv.Atoi(args[1])
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val, err := strconv.Atoi(context.GetMatch(index))
		if err != nil {
			return "<ERROR>"
		}
		return strconv.Itoa((val / bucketSize) * bucketSize)
	})
}

var defaultFunctions = map[string]KeyBuilderFunction{
	"bucket": KeyBuilderFunction(kfBucket),
}
