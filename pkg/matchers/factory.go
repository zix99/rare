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

func ToFactory[T Matcher](f LikeFactory[T]) Factory {
	return &factoryWrapper[T]{f}
}
