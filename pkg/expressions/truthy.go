package expressions

import "strings"

const TruthyVal = "1"
const FalsyVal = ""

func Truthy(s string) bool {
	return strings.TrimSpace(s) != FalsyVal
}

func TruthyStr(is bool) string {
	if is {
		return TruthyVal
	}
	return FalsyVal
}
