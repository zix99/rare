package stdlib

import (
	. "rare/pkg/expressions" //lint:ignore ST1001 Legacy
	"strings"
)

func csvItemEncode(s string) string {
	if strings.ContainsAny(s, "\"\r\n") {
		return "\"" + strings.ReplaceAll(s, "\"", "\"\"") + "\""
	}
	if strings.Contains(s, ",") {
		return "\"" + s + "\""
	}
	return s
}

func kfCsv(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) == 0 {
		return stageLiteral("")
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		var sb strings.Builder
		for i := 0; i < len(args); i++ {
			sb.WriteString(csvItemEncode(args[i](context)))
			if i+1 < len(args) {
				sb.WriteString(",")
			}
		}
		return sb.String()
	}), nil
}
