package expressions

import (
	"strings"
	"unicode"
)

// Space-separated tokenizer that respects escaping,
// quotes, and token symbol {}
func splitTokenizedArguments(s string) []string {
	runes := []rune(s)
	args := make([]string, 0)
	var sb strings.Builder

	tokenDepth := 0
	quoted := false
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if r == '\\' { // something is escaped
			i++
			sb.WriteRune(runes[i])
		} else if r == '"' && !quoted {
			quoted = true
			if tokenDepth > 0 {
				sb.WriteRune('"')
			}
		} else if r == '"' && quoted {
			quoted = false
			if tokenDepth > 0 {
				sb.WriteRune('"')
			} else {
				// Always append, even if empty
				args = append(args, sb.String())
				sb.Reset()
			}
		} else if r == '{' && !quoted {
			tokenDepth++
			sb.WriteRune(r)
		} else if r == '}' && !quoted {
			tokenDepth--
			sb.WriteRune(r)
		} else if unicode.IsSpace(r) && sb.Len() > 0 && tokenDepth == 0 && !quoted {
			args = append(args, sb.String())
			sb.Reset()
		} else if !unicode.IsSpace(r) || quoted || tokenDepth > 0 {
			sb.WriteRune(r)
		}
	}

	if sb.Len() > 0 {
		args = append(args, sb.String())
	}

	return args
}
