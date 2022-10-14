package stdlib

import (
	"fmt"
	"rare/pkg/humanize"
	"strconv"
	"strings"

	. "rare/pkg/expressions" //lint:ignore ST1001 Legacy
)

// {prefix string prefix}
func kfPrefix(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 2 {
		return stageLiteral(ErrorArgCount)
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val := args[0](context)
		contains := args[1](context)

		if strings.HasPrefix(val, contains) {
			return val
		}
		return ""
	})
}

// {suffix string suffix}
func kfSuffix(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 2 {
		return stageLiteral(ErrorArgCount)
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val := args[0](context)
		contains := args[1](context)

		if strings.HasSuffix(val, contains) {
			return val
		}
		return ""
	})
}

func kfUpper(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 1 {
		return stageLiteral(ErrorArgCount)
	}
	return func(context KeyBuilderContext) string {
		return strings.ToUpper(args[0](context))
	}
}

func kfLower(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 1 {
		return stageLiteral(ErrorArgCount)
	}
	return func(context KeyBuilderContext) string {
		return strings.ToLower(args[0](context))
	}
}

// {substr {0} }
func kfSubstr(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 3 {
		return stageLiteral(ErrorArgCount)
	}

	return KeyBuilderStage(func(context KeyBuilderContext) string {
		s := args[0](context)
		lenS := len(s)
		if lenS == 0 {
			return ""
		}

		left, err1 := strconv.Atoi(args[1](context))
		length, err2 := strconv.Atoi(args[2](context))
		if err1 != nil || err2 != nil {
			return ErrorType
		}

		if length < 0 {
			length = 0
		}
		if left < 0 { // negative number wrap-around
			left += lenS
			if left < 0 {
				left = 0
			}
		} else if left > lenS {
			left = lenS
		}

		right := left + length

		if right > lenS {
			right = lenS
		}
		return s[left:right]
	})
}

// {select {0} 1}
func kfSelect(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 2 {
		return stageLiteral(ErrorArgCount)
	}

	return KeyBuilderStage(func(context KeyBuilderContext) string {
		s := args[0](context)
		idx, err := strconv.Atoi(args[1](context))
		if err != nil {
			return ErrorType
		}

		return selectField(s, idx)
	})
}

func selectField(s string, idx int) string {
	currIdx := 0
	wordStart := 0
	inDelim := false
	quoted := false

	for i, c := range s {
		if (quoted && c == '"') || (!quoted && (c == ' ' || c == '\t' || c == '\n' || c == ArraySeparator)) {
			if currIdx == idx {
				return s[wordStart:i]
			}
			inDelim = true
			quoted = false
		} else if c == '"' {
			quoted = !quoted
		} else if inDelim {
			wordStart = i
			currIdx++
			inDelim = false
		}
	}

	if currIdx == idx {
		return s[wordStart:]
	}

	return ""
}

// {format str args...}
// just like fmt.Sprintf
func kfFormat(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) < 1 {
		return stageLiteral(ErrorArgCount)
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		format := args[0](context)

		printArgs := make([]interface{}, len(args)-1)
		for idx, stage := range args[1:] {
			printArgs[idx] = stage(context)
		}

		return fmt.Sprintf(format, printArgs...)
	})
}

func kfHumanizeInt(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 1 {
		return stageLiteral(ErrorArgCount)
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val, err := strconv.Atoi(args[0](context))
		if err != nil {
			return ErrorType
		}
		return humanize.Hi32(val)
	})
}

func kfHumanizeFloat(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 1 {
		return stageLiteral(ErrorArgCount)
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val, err := strconv.ParseFloat(args[0](context), 64)
		if err != nil {
			return ErrorType
		}
		return humanize.Hf(val)
	})
}

func kfBytesize(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) < 1 {
		return stageLiteral(ErrorArgCount)
	}

	precision, err := strconv.Atoi(EvalStageIndexOrDefault(args, 1, "0"))
	if err != nil {
		return stageLiteral(ErrorType)
	}

	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val, err := strconv.ParseUint(args[0](context), 10, 64)
		if err != nil {
			return ErrorType
		}
		return humanize.AlwaysByteSize(val, precision)
	})
}

func kfJoin(delim rune) KeyBuilderFunction {
	return func(args []KeyBuilderStage) KeyBuilderStage {
		if len(args) == 0 {
			return stageLiteral("")
		}
		if len(args) == 1 {
			return args[0]
		}
		return KeyBuilderStage(func(context KeyBuilderContext) string {
			var sb strings.Builder
			sb.WriteString(args[0](context))
			for _, arg := range args[1:] {
				sb.WriteRune(delim)
				sb.WriteString(arg(context))
			}
			return sb.String()
		})
	}
}
