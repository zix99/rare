package stdmath

type Context interface {
	GetMatch(int) float64
	GetKey(string) float64
}

type SimpleContext struct {
	namedVals map[string]float64
}

func (s *SimpleContext) GetMatch(idx int) float64 {
	return 0
}

func (s *SimpleContext) GetKey(k string) float64 {
	return s.namedVals[k]
}
