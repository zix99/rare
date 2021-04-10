package aggregation

import (
	"rare/pkg/stringSplitter"
	"sort"
	"strconv"
)

type SubKeyNamedItem struct {
	Item SubKeyItem
	Name string
}

type SubKeyItem struct {
	count      int64
	submatches []int64
}

type SubKeyCounter struct {
	matches   map[string]*SubKeyItem
	subKeyIdx map[string]int
	subKeys   []string
	errors    uint64
	samples   uint64
}

func NewSubKeyCounter() *SubKeyCounter {
	return &SubKeyCounter{
		matches:   make(map[string]*SubKeyItem),
		subKeyIdx: make(map[string]int),
		subKeys:   make([]string, 0),
	}
}

func (s *SubKeyCounter) Sample(element string) {
	splitter := stringSplitter.Splitter{
		S:     element,
		Delim: "\x00",
	}

	key := splitter.Next()
	subkey := splitter.Next()
	sVal, hasVal := splitter.NextOk()

	if hasVal {
		valNum, err := strconv.ParseInt(sVal, 10, 64)
		if err != nil {
			s.errors++
		} else {
			s.SampleValue(key, subkey, valNum)
		}
	} else {
		s.SampleValue(key, subkey, 1)
	}
}

func (s *SubKeyCounter) SampleValue(key, subkey string, count int64) {
	item := s.getOrCreateKeyItem(key)
	item.count += count

	if subkey != "" {
		subKeyIndex := s.getOrCreateSubkeyIndex(subkey)
		item.submatches[subKeyIndex] += count
	}
}

func (s *SubKeyCounter) getOrCreateKeyItem(key string) *SubKeyItem {
	item := s.matches[key]
	if item == nil {
		item = &SubKeyItem{
			count:      0,
			submatches: make([]int64, len(s.subKeys)),
		}
		s.matches[key] = item
	}
	return item
}

func (s *SubKeyCounter) getOrCreateSubkeyIndex(subkey string) int {
	if idx, ok := s.subKeyIdx[subkey]; !ok {
		s.subKeys = append(s.subKeys, subkey)
		s.subKeyIdx[subkey] = len(s.subKeys) - 1

		for _, item := range s.matches {
			item.submatches = append(item.submatches, 0)
		}

		return len(s.subKeys) - 1
	} else {
		return idx
	}
}

func (s *SubKeyCounter) ParseErrors() uint64 {
	return s.errors
}

func (s *SubKeyCounter) Items() []SubKeyNamedItem {
	ret := make([]SubKeyNamedItem, 0, len(s.matches))
	for key, val := range s.matches {
		ret = append(ret, SubKeyNamedItem{*val, key})
	}
	return ret
}

func (s *SubKeyCounter) ItemsSorted(reverse bool) []SubKeyNamedItem {
	items := s.Items()

	sorter := func(i, j int) bool {
		return items[i].Name < items[j].Name
	}

	if reverse {
		sort.Slice(items, func(i, j int) bool {
			return !sorter(i, j)
		})
	} else {
		sort.Slice(items, sorter)
	}

	return items
}

func (s *SubKeyCounter) SubKeys() []string {
	return s.subKeys
}

func (s *SubKeyItem) Count() int64 {
	return s.count
}

func (s *SubKeyItem) Items() []int64 {
	return s.submatches
}
