package stdmath

import (
	"maps"
	"math"
	"slices"
)

// type Operation rune

type (
	OpFunc  func(left, right float64) float64
	OpUnary func(float64) float64

	OpCode string
)

var orderOfOps = [][]OpCode{
	{"^"},
	{">>", "<<"},
	{"*", "/", "%"},
	{"+", "-"},
	{"&&", "||"},
	{"==", "<=", ">=", ">", "<"},
}

var ops = map[OpCode]OpFunc{
	"+":  func(left, right float64) float64 { return left + right },
	"*":  func(left, right float64) float64 { return left * right },
	"-":  func(left, right float64) float64 { return left - right },
	"/":  func(left, right float64) float64 { return left / right },
	"^":  math.Pow,
	"%":  func(left, right float64) float64 { return float64(int64(left) % int64(right)) },
	"<<": func(left, right float64) float64 { return float64(int64(left) << int64(right)) },
	">>": func(left, right float64) float64 { return float64(int64(left) >> int64(right)) },
	"<":  func(left, right float64) float64 { return conditionalOp(left < right) },
	"<=": func(left, right float64) float64 { return conditionalOp(left <= right) },
	">":  func(left, right float64) float64 { return conditionalOp(left > right) },
	">=": func(left, right float64) float64 { return conditionalOp(left >= right) },
	"==": func(left, right float64) float64 { return conditionalOp(left == right) },

	// todo
	"&&": nil,
	"||": nil,
}

var allOpCodes = slices.Collect(maps.Keys(ops))

var uniOps = map[OpCode]OpUnary{
	"-":   func(f float64) float64 { return -f },
	"abs": math.Abs,

	// todo (more)
	"sin": math.Sin,
	"cos": math.Cos,
	"tan": math.Tan,
}

func isOpAtOrBefore(op0, op1 OpCode) bool {
	for _, opSet := range orderOfOps {
		// if op == op0 { // saw op0 first
		// 	return true
		// }
		// if op == op1 { // saw op1 first
		// 	return false
		// }
		has0, has1 := slices.Contains(opSet, op0), slices.Contains(opSet, op1)
		if has0 || has1 {
			return has0
		}
	}
	panic("op not found")
}

// -1 before, 0 same, 1 after
func opCodeOrder(op0, op1 OpCode) int {
	for _, opSet := range orderOfOps {
		has0, has1 := slices.Contains(opSet, op0), slices.Contains(opSet, op1)
		if has0 && has1 {
			return 0
		}
		if has0 {
			return -1
		}
		if has1 {
			return 1
		}
	}
	panic("op not found")
}

// returns 1/0 based on bool
func conditionalOp(truth bool) float64 {
	if truth {
		return 1.0
	}
	return 0.0
}

func truthy(val float64) bool {
	return val != 0.0
}
