package fuzzy

type FuzzyKey interface {
	Distance(other string, abortAt float32) float32
}
