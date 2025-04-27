package termformat

type StringFormatter func(string) string

func PassthruString(s string) string {
	return s
}

type stringExpressionContext struct {
	s string
}

func (s *stringExpressionContext) GetMatch(idx int) string {
	if idx == 0 {
		return s.s
	}
	return ""
}

func (s *stringExpressionContext) GetKey(key string) string {
	switch key {
	case "val", "value":
		return s.s
	}
	return ""
}

type sep string

func (s sep) GetMatch(idx int) string {
	if idx == 0 {
		return string(s)
	}
	return ""
}

func (s sep) GetKey(key string) string {
	switch key {
	case "val", "value":
		return string(s)
	}
	return ""
}

func StringFromExpression(expr string) (StringFormatter, error) {
	kb, err := expandCompileExpression(expr)
	if err != nil {
		return nil, err
	}

	// ctx := &stringExpressionContext{}
	return func(s string) string {
		// ctx.s = s
		return kb.BuildKey(sep(s))
	}, nil
}
