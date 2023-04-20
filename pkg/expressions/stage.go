package expressions

import (
	"strconv"
	"strings"
)

const (
	ArraySeparator       rune   = '\x00'
	ArraySeparatorString string = string(ArraySeparator)
)

// KeyBuilderFunction defines a helper function at runtime
type KeyBuilderFunction func([]KeyBuilderStage) (KeyBuilderStage, error)

// KeyBuilderStage is a stage within the compiled builder
type KeyBuilderStage func(KeyBuilderContext) string

func stageLiteral(s string) KeyBuilderStage {
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		return s
	})
}

func stageSimpleVariable(s string) KeyBuilderStage {
	index, err := strconv.Atoi(s)
	if err != nil {
		return KeyBuilderStage(func(context KeyBuilderContext) string {
			return context.GetKey(s)
		})
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		return context.GetMatch(index)
	})
}

// make a delim-separated array
func MakeArray(args ...string) string {
	var sb strings.Builder
	for i := 0; i < len(args); i++ {
		if i > 0 {
			sb.WriteRune(ArraySeparator)
		}
		sb.WriteString(args[i])
	}
	return sb.String()
}
