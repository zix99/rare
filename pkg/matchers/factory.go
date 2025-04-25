package matchers

type LikeFactory[T Matcher] interface {
	CreateInstance() T
}

type factoryWrapper[T Matcher] struct {
	matcher LikeFactory[T]
}

func (s *factoryWrapper[T]) CreateInstance() Matcher {
	return s.matcher.CreateInstance()
}

// Maps a factory-like interface to a matcher factory
func ToFactory[T Matcher](f LikeFactory[T]) Factory {
	return &factoryWrapper[T]{f}
}

type noFactoryWrapper struct {
	Matcher
}

func (s *noFactoryWrapper) CreateInstance() Matcher {
	return s
}

// Creates a wrapper factory for a matcher that doesn't need an instance factory
func NoFactory(f Matcher) Factory {
	return &noFactoryWrapper{f}
}
