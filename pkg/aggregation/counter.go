package aggregation

import (
	"rare/pkg/expressions"
	"rare/pkg/stringSplitter"
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
	errors  uint64
	samples uint64
}

func NewCounter() *MatchCounter {
	return &MatchCounter{
		matches: make(map[string]*MatchItem),
	}
}

func (s *MatchCounter) Count() uint64 {
	return s.samples
}

func (s *MatchCounter) GroupCount() int {
	return len(s.matches)
}

func (s *MatchCounter) Sample(element string) {
	splitter := stringSplitter.Splitter{
		S:     element,
		Delim: expressions.ArraySeparatorString,
	}
	key := splitter.Next()
	val, hasVal := splitter.NextOk()

	if hasVal {
		valNum, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			s.errors++
		} else {
			s.SampleValue(key, valNum)
		}
	} else {
		s.SampleValue(key, 1)
	}
}

func (s *MatchCounter) SampleValue(element string, count int64) {
	item := s.matches[element]
	if item == nil {
		item = &MatchItem{
			count: 0,
		}
		s.matches[element] = item
	}
	item.count += count
	s.samples++
}

func (s *MatchCounter) ParseErrors() uint64 {
	return s.errors
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
