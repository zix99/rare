package slicepool

type IntPool struct {
	size int
	pool []int
}

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
