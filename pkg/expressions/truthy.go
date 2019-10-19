package expressions

import "strings"

func Truthy(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	return true
}
