package aggregation

import (
	"sort"
	"strconv"
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

func (s *MatchCounter) Sample(element string) {
	item := s.matches[element]
	if item == nil {
		item = &MatchItem{
			count: 0,
		}
		s.matches[element] = item
	}
	item.count++
}

func (s *MatchCounter) Iter() <-chan MatchPair {
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

	sorter := func(i, j int) bool {
		c0 := items[i].Item.count
		c1 := items[j].Item.count
		if c0 == c1 {
			return items[i].Name < items[j].Name
		}
		return c0 > c1
	}

	if reverse {
		sort.Slice(items, func(i, j int) bool {
			return !sorter(i, j)
		})
	} else {
		sort.Slice(items, sorter)
	}
	return minSlice(items, count)
}

func (s *MatchCounter) ItemsSortedByKey(count int, reverse bool) []MatchPair {
	items := s.Items()

	smartKeySort := func(i, j int) bool {
		num0, err0 := strconv.ParseFloat(items[i].Name, 64)
		num1, err1 := strconv.ParseFloat(items[j].Name, 64)
		if err0 != nil || err1 != nil {
			return items[i].Name < items[j].Name
		}
		return num0 < num1
	}

	if reverse {
		sort.Slice(items, func(i, j int) bool {
			return !smartKeySort(i, j)
		})
	} else {
		sort.Slice(items, smartKeySort)
	}
	return minSlice(items, count)
}

func (s *MatchCounter) ItemsTop(count int) []MatchPair {
	return s.ItemsSorted(count, false)
}

func (s *MatchItem) Count() int64 {
	return s.count
}
