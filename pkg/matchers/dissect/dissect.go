package dissect

import (
	"rare/pkg/slicepool"
	"strings"
	"unsafe"
)

// https://www.elastic.co/guide/en/logstash/current/plugins-filters-dissect.html

// Because of how rare works, and the need to implement `FindSubmatchIndex`
// this is a subset of functionality
// %{key} -- Named key
// %{} or %{?key} -- Named skipped key
// Does NOT support reference keys directly

// Like fastregex, Dissect is NOT thread-safe, and an instance should be created
// per-thread, or it should be locked. This is primarily because of the memory pool

type token struct {
	name, until string
	skip        bool
}

type Dissect struct {
	tokens  []token
	prefix  string
	indexOf func(src, of string) int

	groupNames map[string]int
	groupCount int
}

type DissectInstance struct {
	*Dissect
	groupPool *slicepool.IntPool
}

func CompileEx(expr string, ignoreCase bool) (*Dissect, error) {

	parts := make([]token, 0)
	groupNames := make(map[string]int)
	var prefix string

	groupIndex := 0
	for {
		start := strings.Index(expr, "%{")
		if start < 0 {
			if len(parts) == 0 { // no tokens in expr
				prefix = expr
			}
			break
		}
		if len(parts) == 0 {
			prefix = expr[:start]
		}
		expr = expr[start+2:]

		stop := strings.Index(expr, "}")
		if stop < 0 {
			return nil, ErrorUnclosedToken
		}

		keyName := expr[:stop]
		expr = expr[stop+1:]

		// end is the next token OR end of expr
		end := strings.Index(expr, "%")
		if end < 0 {
			end = len(expr)
		} else if end == 0 {
			return nil, ErrorSequentialToken
		}
		keyUntil := expr[:end]
		expr = expr[end:]

		if ignoreCase {
			keyUntil = strings.ToLower(keyUntil)
		}

		// Special flags
		skipped := false
		switch {
		case len(keyName) == 0: // empty skip
			skipped = true
		case keyName[0] == '?': // named skip
			skipped = true
			keyName = keyName[1:]
		}

		parts = append(parts, token{
			name:  keyName,
			until: keyUntil,
			skip:  skipped,
		})

		if !skipped {
			if _, ok := groupNames[keyName]; ok {
				return nil, ErrorKeyConflict
			}
			groupIndex++
			groupNames[keyName] = groupIndex
		}
	}

	indexOfFunc := strings.Index
	if ignoreCase {
		indexOfFunc = indexIgnoreCase
		prefix = strings.ToLower(prefix)
	}

	return &Dissect{
		groupNames: groupNames,
		groupCount: groupIndex,
		tokens:     parts,
		prefix:     prefix,
		indexOf:    indexOfFunc,
	}, nil
}

func Compile(expr string) (*Dissect, error) {
	return CompileEx(expr, false)
}

func MustCompile(expr string) *Dissect {
	d, err := Compile(expr)
	if err != nil {
		panic(err)
	}
	return d
}

func (s *Dissect) CreateInstance() *DissectInstance {
	return &DissectInstance{
		s,
		slicepool.NewIntPool((s.groupCount*2 + 2) * 1024),
	}
}

// returns indexes of match [first, last, key0Start, key0End, key1Start, ...]
// nil on no match
// replicates logic from regex
func (s *DissectInstance) FindSubmatchIndex(b []byte) []int {
	str := *(*string)(unsafe.Pointer(&b))

	start := 0
	if s.prefix != "" {
		start = s.indexOf(str, s.prefix)
		if start < 0 {
			return nil
		}
		start += len(s.prefix)
	}

	ret := s.groupPool.Get(s.groupCount*2 + 2)
	ret[0] = start - len(s.prefix)

	idx := 2
	for _, token := range s.tokens {

		endOffset := 0
		if token.until == "" {
			endOffset = len(str[start:])
		} else {
			endOffset = s.indexOf(str[start:], token.until)
			if endOffset < 0 {
				return nil
			}
		}

		if !token.skip {
			ret[idx] = start
			ret[idx+1] = start + endOffset
			idx += 2
		}
		start = start + endOffset + len(token.until)
	}

	ret[1] = start

	return ret
}

// Map of key-names to index's in FindSubmatchIndex's return
func (s *Dissect) SubexpNameTable() map[string]int {
	return s.groupNames
}
