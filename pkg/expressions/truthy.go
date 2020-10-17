package expressions

import "strings"

const TruthyVal = "1"
const FalsyVal = ""

func Truthy(s string) bool {
	return strings.TrimSpace(s) != FalsyVal
}
