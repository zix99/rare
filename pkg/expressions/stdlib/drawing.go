package stdlib

import (
	"rare/pkg/color"
	. "rare/pkg/expressions" //lint:ignore ST1001 Legacy
	"rare/pkg/multiterm/termunicode"
	"strconv"
	"strings"
)

func kfColor(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 2 {
		return stageLiteral(ErrorArgCount)
	}
	colorCode, _ := color.LookupColorByName(EvalStageOrDefault(args[0], ""))

	return KeyBuilderStage(func(context KeyBuilderContext) string {
		return color.Wrap(colorCode, args[1](context))
	})
}

// {repeat c {count}}
func kfRepeat(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 2 {
		return stageLiteral(ErrorArgCount)
	}

	char := EvalStageOrDefault(args[0], "|")

	return KeyBuilderStage(func(context KeyBuilderContext) string {
		count, err := strconv.Atoi(args[1](context))
		if err != nil {
			return ErrorType
		}
		return strings.Repeat(char, count)
	})
}

// {bar {val} "maxVal" "len"}
func kfBar(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 3 {
		return stageLiteral(ErrorArgCount)
	}

	maxVal, err := strconv.ParseInt(EvalStageOrDefault(args[1], ""), 10, 64)
	if err != nil {
		return stageLiteral(ErrorType)
	}
	maxLen, err := strconv.ParseInt(EvalStageOrDefault(args[2], ""), 10, 64)
	if err != nil {
		return stageLiteral(ErrorType)
	}

	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val, err := strconv.ParseInt(args[0](context), 10, 64)
		if err != nil {
			return ErrorType
		}
		return termunicode.BarString(val, maxVal, maxLen)
	})
}
