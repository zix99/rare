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
	case "val", "value", ".":
		return s.s
	}
	return ""
}

// Create string->string expression
// Single-thread use only
func StringFromExpression(expr string) (StringFormatter, error) {
	kb, err := expandCompileExpression(expr)
	if err != nil {
		return nil, err
	}

	ctx := &stringExpressionContext{}
	return func(s string) string {
		ctx.s = s
		return kb.BuildKey(ctx)
	}, nil
}
