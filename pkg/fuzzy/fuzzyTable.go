package fuzzy

type fuzzyItem struct {
	key      FuzzyKey
	original string
}

type FuzzyTable struct {
	keys      []fuzzyItem
	matchDist float32
}

func NewFuzzyTable(matchDist float32) *FuzzyTable {
	return &FuzzyTable{
		keys:      make([]fuzzyItem, 0),
		matchDist: matchDist,
	}
}

func (s *FuzzyTable) GetMatchId(val string) (id int, isNew bool) {
	for idx, ele := range s.keys {
		d := ele.key.Distance(val)
		if d > s.matchDist {
			return idx, false
		}
	}

	newItem := fuzzyItem{
		key:      NewLevenshteinKey(val, 0.5),
		original: val,
	}
	s.keys = append(s.keys, newItem)

	return len(s.keys) - 1, true
}

func (s *FuzzyTable) Count() int {
	return len(s.keys)
}

func (s *FuzzyTable) GetString(id int) string {
	return s.keys[id].original
}
