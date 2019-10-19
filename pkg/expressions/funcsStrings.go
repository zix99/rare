package expressions

import (
	"fmt"
	"rare/pkg/humanize"
	"strconv"
	"strings"
)

// {prefix string prefix}
func kfPrefix(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 2 {
		return stageError(ErrorArgCount)
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
		return stageError(ErrorArgCount)
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

// {format str args...}
// just like fmt.Sprintf
func kfFormat(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) < 1 {
		return stageError(ErrorArgCount)
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
		return stageError(ErrorArgCount)
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val, err := strconv.Atoi(args[0](context))
		if err != nil {
			return ErrorType
		}
		return humanize.Hi(val)
	})
}

func kfHumanizeFloat(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 1 {
		return stageError(ErrorArgCount)
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val, err := strconv.ParseFloat(args[0](context), 64)
		if err != nil {
			return ErrorType
		}
		return humanize.Hf(val)
	})
}

var byteSizes = [...]string{"B", "KB", "MB", "GB", "TB", "PB"}

func kfBytesize(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) < 1 {
		return stageError(ErrorArgCount)
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val, err := strconv.Atoi(args[0](context))
		if err != nil {
			return ErrorType
		}

		labelIdx := 0
		for val >= 1024 && labelIdx < len(byteSizes)-1 {
			val = val / 1024
			labelIdx++
		}

		return fmt.Sprintf("%d %s", val, byteSizes[labelIdx])
	})
}
