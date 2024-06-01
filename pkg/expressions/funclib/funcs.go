package funclib

import (
	"rare/pkg/expressions"
	"rare/pkg/expressions/stdlib"
	"rare/pkg/logger"
)

type FunctionSet map[string]expressions.KeyBuilderFunction

var Functions FunctionSet = mapMerge(
	stdlib.StandardFunctions)

func mapMerge[T comparable, Q any](maps ...map[T]Q) (ret map[T]Q) {
	ret = make(map[T]Q)
	for _, m := range maps {
		for k, v := range m {
			ret[k] = v
		}
	}
	return ret
}

func AddFunctions(funcs FunctionSet) {
	for name, fnc := range funcs {
		Functions[name] = fnc
	}
}

func TryAddFunctions(funcs FunctionSet, err error) error {
	if err != nil {
		logger.Printf("Error adding functions: %s", err)
	}
	if funcs != nil {
		AddFunctions(funcs)
	}
	return nil
}
