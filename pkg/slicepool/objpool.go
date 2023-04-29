package slicepool

import "sync"

// A thread-safe pool object that can return or receive points to an object
// Technically can accept objects it didn't create, though that's not good as will pollute the size
// operates in non-blocking mode (it will create a new object if it doesn't have one readily available)
type ObjectPool[T any] struct {
	pool []*T
	m    sync.Mutex
}

// Create an object pool of an initial size. May grow later
func NewObjectPool[T any](size int) *ObjectPool[T] {
	ret := &ObjectPool[T]{
		pool: make([]*T, size),
	}
	for i := 0; i < size; i++ {
		ret.pool[i] = new(T)
	}
	return ret
}

func (s *ObjectPool[T]) Get() (ret *T) {
	s.m.Lock()
	defer s.m.Unlock()

	if len(s.pool) == 0 {
		return new(T)
	}

	end := len(s.pool) - 1
	ret = s.pool[end]
	s.pool = s.pool[:end]
	return
}

func (s *ObjectPool[T]) Return(obj *T) {
	s.m.Lock()
	defer s.m.Unlock()

	s.pool = append(s.pool, obj)
}
