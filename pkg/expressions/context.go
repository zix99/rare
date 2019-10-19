package expressions

// KeyBuilderContext defines how to get information during run-time
type KeyBuilderContext interface {
	GetMatch(idx int) string
}
