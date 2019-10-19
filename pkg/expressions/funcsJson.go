package expressions

import "github.com/tidwall/gjson"

func kfJson(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 2 {
		return stageError(ErrorArgCount)
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		json := args[0](context)
		expression := args[1](context)
		return gjson.Get(json, expression).String()
	})
}
