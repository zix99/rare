package stdlib

import . "rare/pkg/expressions" //lint:ignore ST1001 Legacy

var StandardFunctions = map[string]KeyBuilderFunction{
	"coalesce":  KeyBuilderFunction(kfCoalesce),
	"bucket":    KeyBuilderFunction(kfBucket),
	"clamp":     KeyBuilderFunction(kfClamp),
	"expbucket": KeyBuilderFunction(kfExpBucket),
	"bytesize":  KeyBuilderFunction(kfBytesize),

	// Checks
	"isint": KeyBuilderFunction(kfIsInt),
	"isnum": KeyBuilderFunction(kfIsNum),

	// Arithmatic
	"sumi":  arithmaticHelperi(func(a, b int) int { return a + b }),
	"subi":  arithmaticHelperi(func(a, b int) int { return a - b }),
	"multi": arithmaticHelperi(func(a, b int) int { return a * b }),
	"divi":  arithmaticHelperi(func(a, b int) int { return a / b }),
	"sumf":  arithmaticHelperf(func(a, b float64) float64 { return a + b }),
	"subf":  arithmaticHelperf(func(a, b float64) float64 { return a - b }),
	"multf": arithmaticHelperf(func(a, b float64) float64 { return a * b }),
	"divf":  arithmaticHelperf(func(a, b float64) float64 { return a / b }),

	// Comparisons
	"if": KeyBuilderFunction(kfIf),
	"eq": stringComparator(func(a, b string) string {
		if a == b {
			return a
		}
		return ""
	}),
	"neq": stringComparator(func(a, b string) string {
		if a != b {
			return a
		}
		return ""
	}),
	"not": KeyBuilderFunction(kfNot),
	"lt":  arithmaticEqualityHelper(func(a, b float64) bool { return a < b }),
	"gt":  arithmaticEqualityHelper(func(a, b float64) bool { return a > b }),
	"lte": arithmaticEqualityHelper(func(a, b float64) bool { return a <= b }),
	"gte": arithmaticEqualityHelper(func(a, b float64) bool { return a >= b }),
	"and": KeyBuilderFunction(kfAnd),
	"or":  KeyBuilderFunction(kfOr),

	// Strings
	"like":   KeyBuilderFunction(kfLike),
	"prefix": KeyBuilderFunction(kfPrefix),
	"suffix": KeyBuilderFunction(kfSuffix),
	"format": KeyBuilderFunction(kfFormat),
	"substr": KeyBuilderFunction(kfSubstr),
	"select": KeyBuilderFunction(kfSelect),

	// Separation (Join)
	"tab": kfJoin('\t'),
	"$":   kfJoin(ArraySeparator),

	// Pathing
	"basename": kfPathBase,
	"dirname":  kfPathDir,
	"extname":  kfPathExt,

	// Formatting
	"hi": KeyBuilderFunction(kfHumanizeInt),
	"hf": KeyBuilderFunction(kfHumanizeFloat),

	// Json
	"json": KeyBuilderFunction(kfJsonQuery),

	// CSV
	"csv": KeyBuilderFunction(kfCsv),

	// Time
	"time":           KeyBuilderFunction(kfTimeParse),
	"timeformat":     KeyBuilderFunction(kfTimeFormat),
	"timeattr":       KeyBuilderFunction(kfTimeAttr),
	"buckettime":     KeyBuilderFunction(kfBucketTime),
	"duration":       KeyBuilderFunction(kfDuration),
	"durationformat": KeyBuilderFunction(kfDurationFormat),

	// Color and drawing
	"color":  KeyBuilderFunction(kfColor),
	"repeat": KeyBuilderFunction(kfRepeat),
	"bar":    KeyBuilderFunction(kfBar),
}
