package funclib

import (
	"rare/pkg/expressions"
	"rare/pkg/expressions/stdlib"
	"rare/pkg/logger"
)

type FunctionSet map[string]expressions.KeyBuilderFunction

var Builtins FunctionSet = mapMerge(
	stdlib.StandardFunctions)

var Additional FunctionSet = make(FunctionSet)

func AddFunctions(funcs FunctionSet) {
	for name, fnc := range funcs {
		Additional[name] = fnc
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

func FunctionExists(name string) bool {
	if _, has := Builtins[name]; has {
		return true
	}
	if _, has := Additional[name]; has {
		return true
	}
	return false
}

func mapMerge[T comparable, Q any](maps ...map[T]Q) (ret map[T]Q) {
	ret = make(map[T]Q)
	for _, m := range maps {
		for k, v := range m {
			ret[k] = v
		}
	}
	return ret
}
