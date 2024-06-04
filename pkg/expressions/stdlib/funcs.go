package stdlib

import (
	"math"
	. "rare/pkg/expressions" //lint:ignore ST1001 Legacy
)

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
	"maxi": arithmaticHelperi(func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}),
	"mini": arithmaticHelperi(func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}),
	"sumf":  arithmaticHelperf(func(a, b float64) float64 { return a + b }),
	"subf":  arithmaticHelperf(func(a, b float64) float64 { return a - b }),
	"multf": arithmaticHelperf(func(a, b float64) float64 { return a * b }),
	"divf":  arithmaticHelperf(func(a, b float64) float64 { return a / b }),
	"ceil":  unaryArithmaticHelperfi(func(f float64) int64 { return int64(math.Ceil(f)) }),
	"floor": unaryArithmaticHelperfi(func(f float64) int64 { return int64(math.Floor(f)) }),
	"round": kfRound,

	// Comparisons
	"if":     KeyBuilderFunction(kfIf),
	"switch": kfSwitch,
	"unless": KeyBuilderFunction(kfUnless),
	"eq": stringComparator(func(a, b string) string {
		if a == b {
			return TruthyVal
		}
		return FalsyVal
	}),
	"neq": stringComparator(func(a, b string) string {
		if a != b {
			return TruthyVal
		}
		return FalsyVal
	}),
	"not": KeyBuilderFunction(kfNot),
	"lt":  arithmaticEqualityHelper(func(a, b float64) bool { return a < b }),
	"gt":  arithmaticEqualityHelper(func(a, b float64) bool { return a > b }),
	"lte": arithmaticEqualityHelper(func(a, b float64) bool { return a <= b }),
	"gte": arithmaticEqualityHelper(func(a, b float64) bool { return a >= b }),
	"and": KeyBuilderFunction(kfAnd),
	"or":  KeyBuilderFunction(kfOr),

	// Strings
	"len":    KeyBuilderFunction(kfLen),
	"like":   KeyBuilderFunction(kfLike),
	"prefix": KeyBuilderFunction(kfPrefix),
	"suffix": KeyBuilderFunction(kfSuffix),
	"format": KeyBuilderFunction(kfFormat),
	"substr": KeyBuilderFunction(kfSubstr),
	"select": KeyBuilderFunction(kfSelect),
	"upper":  KeyBuilderFunction(kfUpper),
	"lower":  KeyBuilderFunction(kfLower),

	// Separation (Join)
	"tab": kfJoin('\t'),
	"$":   kfJoin(ArraySeparator),

	// Ranges
	"@":       kfJoin(ArraySeparator),
	"@len":    kfArrayLen,
	"@map":    kfArrayMap,
	"@split":  kfArraySplit,
	"@select": kfArraySelect,
	"@join":   kfArrayJoin,
	"@reduce": kfArrayReduce,
	"@filter": kfArrayFilter,
	"@slice":  kfArraySlice,
	"@in":     kfArrayIn,

	// Pathing
	"basename": kfPathBase,
	"dirname":  kfPathDir,
	"extname":  kfPathExt,

	// File operations
	"load":   kfLoadFile,
	"lookup": kfLookupKey,
	"haskey": kfHasKey,

	// Formatting
	"hi":      KeyBuilderFunction(kfHumanizeInt),
	"hf":      KeyBuilderFunction(kfHumanizeFloat),
	"percent": kfPercent,

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
