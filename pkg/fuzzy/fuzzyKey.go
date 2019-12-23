package fuzzy

type FuzzyKey interface {
	Distance(other string) float32
}
