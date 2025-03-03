package stdlib

import (
	"fmt"
	"rare/pkg/humanize"
	"strconv"
	"strings"

	. "rare/pkg/expressions" //lint:ignore ST1001 Legacy
)

// {len string}
func kfLen(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 1 {
		return stageErrArgCount(args, 1)
	}
	return func(context KeyBuilderContext) string {
		return strconv.Itoa(len(args[0](context)))
	}, nil
}

// {prefix string prefix}
func kfPrefix(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 2 {
		return stageErrArgCount(args, 2)
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val := args[0](context)
		contains := args[1](context)

		if strings.HasPrefix(val, contains) {
			return val
		}
		return ""
	}), nil
}

// {suffix string suffix}
func kfSuffix(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 2 {
		return stageErrArgCount(args, 2)
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val := args[0](context)
		contains := args[1](context)

		if strings.HasSuffix(val, contains) {
			return val
		}
		return ""
	}), nil
}

func kfUpper(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 1 {
		return stageErrArgCount(args, 1)
	}
	return func(context KeyBuilderContext) string {
		return strings.ToUpper(args[0](context))
	}, nil
}

func kfLower(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 1 {
		return stageErrArgCount(args, 1)
	}
	return func(context KeyBuilderContext) string {
		return strings.ToLower(args[0](context))
	}, nil
}

// {substr {0} }
func kfSubstr(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 3 {
		return stageErrArgCount(args, 3)
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
			return ErrorNum
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
	}), nil
}

// {select {0} 1}
func kfSelect(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 2 {
		return stageErrArgCount(args, 2)
	}

	return KeyBuilderStage(func(context KeyBuilderContext) string {
		s := args[0](context)
		idx, err := strconv.Atoi(args[1](context))
		if err != nil {
			return ErrorNum
		}

		return selectField(s, idx)
	}), nil
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
func kfFormat(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) < 1 {
		return stageErrArgRange(args, "1+")
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		format := args[0](context)

		printArgs := make([]interface{}, len(args)-1)
		for idx, stage := range args[1:] {
			printArgs[idx] = stage(context)
		}

		return fmt.Sprintf(format, printArgs...)
	}), nil
}

// {percent val [decimals=1] [[min] max]}
func kfPercent(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if !isArgCountBetween(args, 1, 4) {
		return stageErrArgRange(args, "1-4")
	}

	decimals, hasDecimals := EvalArgInt(args, 1, 1)
	if !hasDecimals {
		return stageArgError(ErrConst, 1)
	}

	switch len(args) {
	case 3: // max, no min (0)
		return func(context KeyBuilderContext) string {
			max, err := strconv.ParseFloat(args[2](context), 64)
			if err != nil {
				return ErrorNum
			}

			val, err := strconv.ParseFloat(args[0](context), 64)
			if err != nil {
				return ErrorNum
			}
			return strconv.FormatFloat(val*100.0/max, 'f', decimals, 64) + "%"
		}, nil
	case 4: // min, max
		return func(context KeyBuilderContext) string {
			min, err := strconv.ParseFloat(args[2](context), 64)
			if err != nil {
				return ErrorNum
			}

			max, err := strconv.ParseFloat(args[3](context), 64)
			if err != nil {
				return ErrorNum
			}

			val, err := strconv.ParseFloat(args[0](context), 64)
			if err != nil {
				return ErrorNum
			}
			return strconv.FormatFloat((val-min)*100.0/(max-min), 'f', decimals, 64) + "%"
		}, nil
	default:
		return func(context KeyBuilderContext) string {
			val, err := strconv.ParseFloat(args[0](context), 64)
			if err != nil {
				return ErrorNum
			}
			return strconv.FormatFloat(val*100.0, 'f', decimals, 64) + "%"
		}, nil
	}
}

func kfHumanizeInt(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 1 {
		return stageErrArgCount(args, 1)
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val, err := strconv.Atoi(args[0](context))
		if err != nil {
			return ErrorNum
		}
		return humanize.Hi32(val)
	}), nil
}

func kfHumanizeFloat(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 1 {
		return stageErrArgCount(args, 1)
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val, err := strconv.ParseFloat(args[0](context), 64)
		if err != nil {
			return ErrorNum
		}
		return humanize.Hf(val)
	}), nil
}

// {bytesize val [precision]}
func kfBytesize(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if !isArgCountBetween(args, 1, 2) {
		return stageErrArgRange(args, "1-2")
	}

	precision, pOk := EvalArgInt(args, 1, 0)
	if !pOk {
		return stageArgError(ErrNum, 1)
	}

	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val, err := strconv.ParseUint(args[0](context), 10, 64)
		if err != nil {
			return ErrorNum
		}
		return humanize.AlwaysByteSize(val, precision)
	}), nil
}

// {bytesizesi val [precision]}
func kfBytesizeSi(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if !isArgCountBetween(args, 1, 2) {
		return stageErrArgRange(args, "1-2")
	}

	precision, pOk := EvalArgInt(args, 1, 0)
	if !pOk {
		return stageArgError(ErrNum, 1)
	}

	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val, err := strconv.ParseUint(args[0](context), 10, 64)
		if err != nil {
			return ErrorNum
		}
		return humanize.AlwaysByteSizeSi(val, precision)
	}), nil
}

// {downscale val [precision]}
func kfDownscale(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if !isArgCountBetween(args, 1, 2) {
		return stageErrArgRange(args, "1-2")
	}

	precision, pOk := EvalArgInt(args, 1, 0)
	if !pOk {
		return stageArgError(ErrNum, 1)
	}

	return func(context KeyBuilderContext) string {
		val, err := strconv.ParseInt(args[0](context), 10, 64)
		if err != nil {
			return ErrorNum
		}
		return humanize.AlwaysDownscale(val, precision)
	}, nil
}

func kfJoin(delim rune) KeyBuilderFunction {
	return func(args []KeyBuilderStage) (KeyBuilderStage, error) {
		if len(args) == 0 {
			return stageLiteral("")
		}
		if len(args) == 1 {
			return args[0], nil
		}
		return KeyBuilderStage(func(context KeyBuilderContext) string {
			var sb strings.Builder
			sb.WriteString(args[0](context))
			for _, arg := range args[1:] {
				sb.WriteRune(delim)
				sb.WriteString(arg(context))
			}
			return sb.String()
		}), nil
	}
}
