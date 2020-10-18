package stdlib

import (
	. "rare/pkg/expressions" //lint:ignore ST1001 Legacy
	"strconv"
)

func kfIsInt(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 1 {
		return stageLiteral(ErrorArgCount)
	}

	return func(context KeyBuilderContext) string {
		_, err := strconv.Atoi(args[0](context))
		if err != nil {
			return FalsyVal
		}
		return TruthyVal
	}
}

func kfIsNum(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 1 {
		return stageLiteral(ErrorArgCount)
	}

	return func(context KeyBuilderContext) string {
		_, err := strconv.ParseFloat(args[0](context), 64)
		if err != nil {
			return FalsyVal
		}
		return TruthyVal
	}
}
