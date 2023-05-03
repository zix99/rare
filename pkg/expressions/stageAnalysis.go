package expressions

import (
	"strconv"
)

// monitorContext allows monitoring of context use
//   largely for static analysis of an expression
type monitorContext struct {
	keyLookups int
}

func (s *monitorContext) GetMatch(idx int) string {
	s.keyLookups++
	return ""
}

func (s *monitorContext) GetKey(key string) string {
	s.keyLookups++
	return ""
}

func EvalStaticStage(stage KeyBuilderStage) (ret string, ok bool) {
	var monitor monitorContext
	ret = stage(&monitor)
	ok = (monitor.keyLookups == 0)
	return
}

// Eval a stage, but upon any error, returns default instead
// Deprecated: Consider using EvalStaticStage for better error detection
func EvalStageIndexOrDefault(stages []KeyBuilderStage, idx int, dflt string) string {
	if idx < len(stages) {
		if val, ok := EvalStaticStage(stages[idx]); ok {
			return val
		}
	}
	return dflt
}

// Evals stage and parses as int. If fails for any reason (stage eval or parsing), will return false
func EvalStageInt(stage KeyBuilderStage) (int, bool) {
	if s, ok := EvalStaticStage(stage); ok {
		if v, err := strconv.Atoi(s); err == nil {
			return v, true
		}
	}
	return 0, false
}

// Evals stage and parses as int. If fails for any reason (stage eval or parsing), will return false
func EvalStageInt64(stage KeyBuilderStage) (int64, bool) {
	if s, ok := EvalStaticStage(stage); ok {
		if v, err := strconv.ParseInt(s, 10, 64); err == nil {
			return v, true
		}
	}
	return 0, false
}

// Helper to EvalStageInt
func EvalArgInt(stages []KeyBuilderStage, idx int, dflt int) (int, bool) {
	if idx < len(stages) {
		return EvalStageInt(stages[idx])
	}
	return dflt, true
}
