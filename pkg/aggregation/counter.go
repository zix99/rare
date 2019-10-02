package aggregation

type MatchItem struct {
	count int64
}

type MatchPair struct {
	Item MatchItem
	Name string
}

type MatchCounter struct {
	matches map[string]*MatchItem
}

func New() *MatchCounter {
	return &MatchCounter{
		matches: make(map[string]*MatchItem),
	}
}

func (s *MatchCounter) Inc(element string) {
	item := s.matches[element]
	if item == nil {
		item = &MatchItem{
			count: 0,
		}
		s.matches[element] = item
	}
	item.count++
}

func (s *MatchCounter) Iter() chan MatchPair {
	c := make(chan MatchPair)
	go func() {
		for key, value := range s.matches {
			select {
			case c <- MatchPair{
				Item: *value,
				Name: key,
			}:
			case <-c:
				close(c)
				return
			}
			close(c)
		}
	}()
	return c
}
