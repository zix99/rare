package stdlib

import (
	"github.com/tidwall/gjson"

	. "rare/pkg/expressions" //lint:ignore ST1001 Legacy
)

func kfJsonQuery(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) == 1 {
		// Assumes "{0}" is the json blob to extract, so arg[0] is the key
		return KeyBuilderStage(func(context KeyBuilderContext) string {
			json := context.GetMatch(0)
			expression := args[0](context)
			return gjson.Get(json, expression).String()
		})
	} else if len(args) == 2 {
		// Json is arg[0], key is arg[1]
		return KeyBuilderStage(func(context KeyBuilderContext) string {
			json := args[0](context)
			expression := args[1](context)
			return gjson.Get(json, expression).String()
		})
	} else {
		return stageLiteral(ErrorArgCount)
	}
}
