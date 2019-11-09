package expressions

import "strings"

const TruthyVal = "1"
const FalsyVal = ""

func Truthy(s string) bool {
	s = strings.TrimSpace(s)
	if s == FalsyVal {
		return false
	}
	return true
}
