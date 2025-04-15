package expressions

// monitorContext allows monitoring of context use
//
//	largely for static analysis of an expression
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
