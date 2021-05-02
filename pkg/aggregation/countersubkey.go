package aggregation

import (
	"rare/pkg/stringSplitter"
	"sort"
	"strconv"
)

// SubKeyNamedItem is a returnable key-pair for data aggregation
type SubKeyNamedItem struct {
	Item SubKeyItem
	Name string
}

type SubKeyItem struct {
	count      int64
	submatches []int64 // Matches, in same order as subKeys
}

type SubKeyCounter struct {
	matches   map[string]*SubKeyItem // All matches of the top level key
	subKeys   []string               // All the subkeys. Never changes order (append-only)
	subKeyIdx map[string]int         // Name to index of subKeys
	errors    uint64
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

	subKeyIndex := s.getOrCreateSubkeyIndex(subkey)
	item.submatches[subKeyIndex] += count
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
		s.subKeys, idx = insertAlphanumeric(s.subKeys, subkey)
		s.subKeyIdx[subkey] = idx

		for _, item := range s.matches {
			item.submatches = insertAti64(item.submatches, idx, 0)
		}

		return idx
	} else {
		return idx
	}
}

// insertAlphanumeric assumes slice is already in-order, and inserts a new element
//  and returns the index it was inserted at
func insertAlphanumeric(slice []string, ele string) (ret []string, idx int) {
	for i, val := range slice {
		if ele < val {
			ret = insertAt(slice, i, ele)
			idx = i
			return
		}
	}

	// Must be at the end, it's a simple case
	idx = len(slice)
	ret = append(slice, ele)
	return
}

func insertAt(slice []string, idx int, ele string) (ret []string) {
	ret = append(slice, "")
	copy(ret[idx+1:], ret[idx:])
	ret[idx] = ele
	return
}

func insertAti64(slice []int64, idx int, ele int64) (ret []int64) {
	ret = append(slice, 0)
	copy(ret[idx+1:], ret[idx:])
	ret[idx] = ele
	return
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
