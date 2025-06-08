package stdlib

import (
	"strconv"

	. "github.com/zix99/rare/pkg/expressions" //lint:ignore ST1001 Legacy
)

func kfIsInt(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 1 {
		return stageErrArgCount(args, 1)
	}

	return func(context KeyBuilderContext) string {
		_, err := strconv.Atoi(args[0](context))
		if err != nil {
			return FalsyVal
		}
		return TruthyVal
	}, nil
}

func kfIsNum(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 1 {
		return stageErrArgCount(args, 1)
	}

	return func(context KeyBuilderContext) string {
		_, err := strconv.ParseFloat(args[0](context), 64)
		if err != nil {
			return FalsyVal
		}
		return TruthyVal
	}, nil
}
