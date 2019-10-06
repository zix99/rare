package extractor

import "strconv"

// KeyBuilderContext defines how to get information during run-time
type KeyBuilderContext interface {
	GetMatch(idx int) string
}

// KeyBuilderStage is a stage within the compiled builder
type KeyBuilderStage func(KeyBuilderContext) string

// KeyBuilderFunction defines a helper function at runtime
type KeyBuilderFunction func([]string) KeyBuilderStage

func stageLiteral(s string) KeyBuilderStage {
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		return s
	})
}

func stageSimpleVariable(s string) KeyBuilderStage {
	index, _ := strconv.Atoi(s)
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		return context.GetMatch(index)
	})
}
