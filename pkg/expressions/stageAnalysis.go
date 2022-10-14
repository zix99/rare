package expressions

import "strconv"

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
func EvalStageOrDefault(stage KeyBuilderStage, dflt string) string {
	if val, ok := EvalStaticStage(stage); ok {
		return val
	}
	return dflt
}

func EvalStageIndexOrDefault(stages []KeyBuilderStage, idx int, dflt string) string {
	if idx < len(stages) {
		return EvalStageOrDefault(stages[idx], dflt)
	}
	return dflt
}

func EvalStageInt(stage KeyBuilderStage, dflt int) int {
	if s, ok := EvalStaticStage(stage); ok {
		if v, err := strconv.Atoi(s); err == nil {
			return v
		}
	}
	return dflt
}
