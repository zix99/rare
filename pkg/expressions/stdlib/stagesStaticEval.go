package stdlib

import (
	"rare/pkg/expressions"
	"strconv"
)

// Eval a stage, but upon any error, returns default instead
// Deprecated: Consider using EvalStaticStage for better error detection
func EvalStageIndexOrDefault(stages []expressions.KeyBuilderStage, idx int, dflt string) string {
	if idx < len(stages) {
		if val, ok := expressions.EvalStaticStage(stages[idx]); ok {
			return val
		}
	}
	return dflt
}

// Evals stage and parses as int. If fails for any reason (stage eval or parsing), will return false
func EvalStageInt(stage expressions.KeyBuilderStage) (int, bool) {
	if s, ok := expressions.EvalStaticStage(stage); ok {
		if v, err := strconv.Atoi(s); err == nil {
			return v, true
		}
	}
	return 0, false
}

// Evals stage and parses as int. If fails for any reason (stage eval or parsing), will return false
func EvalStageInt64(stage expressions.KeyBuilderStage) (int64, bool) {
	if s, ok := expressions.EvalStaticStage(stage); ok {
		if v, err := strconv.ParseInt(s, 10, 64); err == nil {
			return v, true
		}
	}
	return 0, false
}

// Helper to EvalStageInt
func EvalArgInt(stages []expressions.KeyBuilderStage, idx int, dflt int) (int, bool) {
	if idx < len(stages) {
		return EvalStageInt(stages[idx])
	}
	return dflt, true
}
