package stdlib

import "rare/pkg/expressions"

func NewStdKeyBuilderEx(autoOptimize bool) *expressions.KeyBuilder {
	kb := expressions.NewKeyBuilderEx(autoOptimize)
	kb.Funcs(StandardFunctions)
	return kb
}

func NewStdKeyBuilder() *expressions.KeyBuilder {
	return NewStdKeyBuilderEx(true)
}
