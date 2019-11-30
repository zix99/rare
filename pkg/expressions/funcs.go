package expressions

// KeyBuilderFunction defines a helper function at runtime
type KeyBuilderFunction func([]KeyBuilderStage) KeyBuilderStage

var defaultFunctions = map[string]KeyBuilderFunction{
	"coalesce":  KeyBuilderFunction(kfCoalesce),
	"bucket":    KeyBuilderFunction(kfBucket),
	"expbucket": KeyBuilderFunction(kfExpBucket),
	"bytesize":  KeyBuilderFunction(kfBytesize),

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
	"tab":    KeyBuilderFunction(kfTab),

	// Formatting
	"hi": KeyBuilderFunction(kfHumanizeInt),
	"hf": KeyBuilderFunction(kfHumanizeFloat),

	// Json
	"json": KeyBuilderFunction(kfJson),

	// Time
	"time":       KeyBuilderFunction(kfTimeParse),
	"timeformat": KeyBuilderFunction(kfTimeFormat),
	"duration":   KeyBuilderFunction(kfDuration),
}
