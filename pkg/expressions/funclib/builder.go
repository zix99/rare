package funclib

import (
	"github.com/zix99/rare/pkg/expressions"
)

func NewKeyBuilderEx(autoOptimize bool) *expressions.KeyBuilder {
	kb := expressions.NewKeyBuilderEx(autoOptimize)
	kb.Funcs(Builtins)
	kb.Funcs(Additional)
	return kb
}

func NewKeyBuilder() *expressions.KeyBuilder {
	return NewKeyBuilderEx(true)
}
