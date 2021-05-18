package testutil

import "regexp"

var stringSplitter = regexp.MustCompile(`([^\s"]+)|"([^"]*)"`)

func SplitQuotedString(s string) []string {
	matches := stringSplitter.FindAllStringSubmatch(s, -1)

	ret := make([]string, 0)
	for _, v := range matches {
		if v[2] != "" {
			ret = append(ret, v[2])
		} else {
			ret = append(ret, v[1])
		}
	}

	return ret
}
