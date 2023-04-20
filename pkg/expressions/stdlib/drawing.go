package stdlib

import (
	"rare/pkg/color"
	. "rare/pkg/expressions" //lint:ignore ST1001 Legacy
	"rare/pkg/multiterm/termunicode"
	"strconv"
	"strings"
)

func kfColor(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 2 {
		return stageError(ErrArgCount)
	}
	colorCode, _ := color.LookupColorByName(EvalStageOrDefault(args[0], ""))

	return KeyBuilderStage(func(context KeyBuilderContext) string {
		return color.Wrap(colorCode, args[1](context))
	}), nil
}

// {repeat c {count}}
func kfRepeat(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 2 {
		return stageError(ErrArgCount)
	}

	char := EvalStageOrDefault(args[0], "|")

	return KeyBuilderStage(func(context KeyBuilderContext) string {
		count, err := strconv.Atoi(args[1](context))
		if err != nil {
			return ErrorNum
		}
		return strings.Repeat(char, count)
	}), nil
}

// {bar {val} "maxVal" "len"}
func kfBar(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 3 {
		return stageError(ErrArgCount)
	}

	maxVal, err := strconv.ParseInt(EvalStageOrDefault(args[1], ""), 10, 64)
	if err != nil {
		return stageError(ErrNum)
	}
	maxLen, err := strconv.ParseInt(EvalStageOrDefault(args[2], ""), 10, 64)
	if err != nil {
		return stageError(ErrNum)
	}

	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val, err := strconv.ParseInt(args[0](context), 10, 64)
		if err != nil {
			return ErrorNum
		}
		return termunicode.BarString(val, maxVal, maxLen)
	}), nil
}
