package aggregation

// Aggregator provides default agg interface
type Aggregator interface {
	Sample(element string)
	ParseErrors() uint64
}
