package aggregation

import (
	"sort"
)

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

func NewCounter() *MatchCounter {
	return &MatchCounter{
		matches: make(map[string]*MatchItem),
	}
}

func (s *MatchCounter) GroupCount() int {
	return len(s.matches)
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
		}
		close(c)
	}()
	return c
}

func (s *MatchCounter) Items() []MatchPair {
	items := make([]MatchPair, 0, len(s.matches))
	for key, val := range s.matches {
		items = append(items, MatchPair{
			Item: *val,
			Name: key,
		})
	}
	return items
}

func minSlice(items []MatchPair, count int) []MatchPair {
	if len(items) < count {
		return items
	}
	return items[:count]
}

func (s *MatchCounter) ItemsSorted(count int, reverse bool) []MatchPair {
	items := s.Items()
	if reverse {
		sort.Slice(items, func(i, j int) bool {
			return items[i].Item.count < items[j].Item.count
		})
	} else {
		sort.Slice(items, func(i, j int) bool {
			return items[i].Item.count > items[j].Item.count
		})
	}
	return minSlice(items, count)
}

func (s *MatchCounter) ItemsSortedByKey(count int, reverse bool) []MatchPair {
	items := s.Items()
	if reverse {
		sort.Slice(items, func(i, j int) bool {
			return items[i].Name > items[j].Name
		})
	} else {
		sort.Slice(items, func(i, j int) bool {
			return items[i].Name < items[j].Name
		})
	}
	return minSlice(items, count)
}

func (s *MatchCounter) ItemsTop(count int) []MatchPair {
	return s.ItemsSorted(count, false)
}

func (s *MatchItem) Count() int64 {
	return s.count
}
