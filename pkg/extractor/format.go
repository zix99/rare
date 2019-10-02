package extractor

import (
	"fmt"
	"strings"
)

// Build a string from a set of matches and a format string like
// "$1 bla bla $2"
func buildStringFromGroups(matches []string, format string) string {
	if len(format) == 0 {
		return strings.Join(matches, ",")
	}

	var temp string = format
	for idx, val := range matches {
		temp = strings.Replace(temp, fmt.Sprintf("$%d", idx), val, -1)
	}
	return temp
}
