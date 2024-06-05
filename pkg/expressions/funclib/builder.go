package funclib

import "rare/pkg/expressions"

func NewKeyBuilderEx(autoOptimize bool) *expressions.KeyBuilder {
	kb := expressions.NewKeyBuilderEx(autoOptimize)
	kb.Funcs(Functions)
	return kb
}

func NewKeyBuilder() *expressions.KeyBuilder {
	return NewKeyBuilderEx(true)
}
