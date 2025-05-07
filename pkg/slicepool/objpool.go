package slicepool

import "sync"

// A thread-safe pool object that can return or receive points to an object
// Technically can accept objects it didn't create, though that's not good as will pollute the size
// operates in non-blocking mode (it will create a new object if it doesn't have one readily available)
//
// Uses sync.pool as underlying store, as its about 2x as fast as mutex lock/unlock, and shaves about
// 25% cpu off of parallelized execution
type ObjectPool[T any] struct {
	pool sync.Pool
}

// Create an object pool of an initial size. May grow later
func NewObjectPool[T any](size int) *ObjectPool[T] {
	return NewObjectPoolEx(size, func() *T { return new(T) })
}

// Create an object pool with a custom object initializer
func NewObjectPoolEx[T any](size int, newer func() *T) *ObjectPool[T] {
	ret := &ObjectPool[T]{}

	ret.pool.New = func() any { return newer() }

	for range size {
		ret.pool.Put(newer())
	}

	return ret
}

func (s *ObjectPool[T]) Get() (ret *T) {
	return s.pool.Get().(*T)
}

func (s *ObjectPool[T]) Return(obj *T) {
	s.pool.Put(obj)
}
