package slicepool

type IntPool struct {
	size int
	pool []int
}

// Create an slice of integers, and release chunks at a team for read-write.
// If run out of pool space, will create a new pool of `size`
// Use when need lots of small int slices, as will limit total allocs
func NewIntPool(size int) *IntPool {
	return &IntPool{
		size: size,
		pool: make([]int, size),
	}
}

func (s *IntPool) Get(n int) (ret []int) {
	if len(s.pool) < n {
		if n > s.size {
			panic("pool not large enough")
		}
		s.pool = make([]int, s.size)
	}
	ret = s.pool[:n]
	s.pool = s.pool[n:]
	return
}
