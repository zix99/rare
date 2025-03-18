package aggregation

import (
	"rare/pkg/aggregation/sorting"
	"rare/pkg/expressions"
	"rare/pkg/stringSplitter"
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
	total   int64
}

func NewCounter() *MatchCounter {
	return &MatchCounter{
		matches: make(map[string]*MatchItem),
	}
}

func (s *MatchCounter) Total() int64 {
	return s.total
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
	s.total += count
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

func (s *MatchCounter) ItemsSortedBy(count int, sorter sorting.NameValueSorter) []MatchPair {
	items := s.Items()
	sorting.SortBy(items, sorter, func(obj MatchPair) sorting.NameValuePair {
		return sorting.NameValuePair{
			Name:  obj.Name,
			Value: obj.Item.count,
		}
	})
	return minSlice(items, count)
}

func (s *MatchItem) Count() int64 {
	return s.count
}
