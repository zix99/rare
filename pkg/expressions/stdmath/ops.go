package stdmath

import "math"

// type Operation rune

type (
	OpFunc  func(left, right float64) float64
	OpUnary func(float64) float64
)

var orderOfOps = []string{"^", "*", "/", "+", "-", "&&", "||", "<=", ">=", ">", "<"}

var ops = map[string]OpFunc{
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

var uniOps = map[string]OpUnary{
	"-":   func(f float64) float64 { return -f },
	"abs": math.Abs,

	// todo
	"sin": math.Sin,
	"cos": math.Cos,
	"tan": math.Tan,
}

func isOpBefore(op0, op1 string) bool {
	for _, op := range orderOfOps {
		if op == op0 { // saw op0 first
			return true
		}
		if op == op1 { // saw op1 first
			return false
		}
	}
	panic("fixme")
}

func conditionalOp(truth bool) float64 {
	if truth {
		return 1.0
	}
	return 0.0
}
