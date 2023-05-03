package stdlib

import (
	"rare/pkg/color"
	. "rare/pkg/expressions" //lint:ignore ST1001 Legacy
	"rare/pkg/multiterm/termunicode"
	"strconv"
	"strings"
)

// {color "color" content}
func kfColor(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 2 {
		return stageErrArgCount(args, 2)
	}
	colorName, colorNameOk := EvalStaticStage(args[0])
	if !colorNameOk {
		return stageArgError(ErrConst, 0)
	}

	colorCode, hasColor := color.LookupColorByName(colorName)
	if !hasColor {
		return stageArgError(ErrEnum, 0)
	}

	return KeyBuilderStage(func(context KeyBuilderContext) string {
		return color.Wrap(colorCode, args[1](context))
	}), nil
}

// {repeat c {count}}
func kfRepeat(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 2 {
		return stageErrArgCount(args, 2)
	}

	char, charOk := EvalStaticStage(args[0])
	if !charOk {
		return stageArgError(ErrConst, 0)
	}

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
		return stageErrArgCount(args, 3)
	}

	maxVal, maxValOk := EvalStageInt64(args[1])
	if !maxValOk {
		return stageArgError(ErrNum, 1)
	}
	maxLen, maxLenOk := EvalStageInt64(args[2])
	if !maxLenOk {
		return stageArgError(ErrNum, 2)
	}

	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val, err := strconv.ParseInt(args[0](context), 10, 64)
		if err != nil {
			return ErrorNum
		}
		return termunicode.BarString(val, maxVal, maxLen)
	}), nil
}
