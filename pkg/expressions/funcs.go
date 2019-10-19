package expressions

var defaultFunctions = map[string]KeyBuilderFunction{
	"coalesce":  KeyBuilderFunction(kfCoalesce),
	"bucket":    KeyBuilderFunction(kfBucket),
	"expbucket": KeyBuilderFunction(kfExpBucket),
	"bytesize":  KeyBuilderFunction(kfBytesize),
	"sumi":      arithmaticHelperi(func(a, b int) int { return a + b }),
	"subi":      arithmaticHelperi(func(a, b int) int { return a - b }),
	"multi":     arithmaticHelperi(func(a, b int) int { return a * b }),
	"divi":      arithmaticHelperi(func(a, b int) int { return a / b }),
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
	"not":    KeyBuilderFunction(kfNot),
	"lt":     arithmaticEqualityHelper(func(a, b int) bool { return a < b }),
	"gt":     arithmaticEqualityHelper(func(a, b int) bool { return a > b }),
	"lte":    arithmaticEqualityHelper(func(a, b int) bool { return a <= b }),
	"gte":    arithmaticEqualityHelper(func(a, b int) bool { return a >= b }),
	"and":    KeyBuilderFunction(kfAnd),
	"or":     KeyBuilderFunction(kfOr),
	"like":   KeyBuilderFunction(kfLike),
	"prefix": KeyBuilderFunction(kfPrefix),
	"suffix": KeyBuilderFunction(kfSuffix),
	"format": KeyBuilderFunction(kfFormat),
	"hi":     KeyBuilderFunction(kfHumanizeInt),
	"hf":     KeyBuilderFunction(kfHumanizeFloat),
	"json":   KeyBuilderFunction(kfJson),
}
