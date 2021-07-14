package minijson

import (
	"strconv"
	"strings"
)

/*
A JSON writer that is explicit and doesn't marshal/unmarshal
*/

type JsonObjectBuilder struct {
	sb       strings.Builder
	keyCount int
}

func (s *JsonObjectBuilder) OpenEx(hint int) {
	if hint > 0 {
		s.sb.Grow(hint)
	}
	s.sb.WriteRune('{')
}

func (s *JsonObjectBuilder) Open() {
	s.OpenEx(0)
}

func (s *JsonObjectBuilder) Close() {
	s.sb.WriteRune('}')
}

func (s *JsonObjectBuilder) String() string {
	return s.sb.String()
}

func (s *JsonObjectBuilder) KeyCount() int {
	return s.keyCount
}

func (s *JsonObjectBuilder) WriteInferred(key, val string) {
	if isNumeric(val) {
		s.WriteLiteral(key, val)
	} else if strings.EqualFold(val, "true") {
		s.WriteLiteral(key, "true")
	} else if strings.EqualFold(val, "false") {
		s.WriteLiteral(key, "false")
	} else {
		s.WriteString(key, val)
	}
}

// Write a {"Key": literal} (Note, no quotes in literal)
func (s *JsonObjectBuilder) WriteLiteral(key, literal string) {
	s.writeKey(key)
	s.sb.WriteString(literal)
}

func (s *JsonObjectBuilder) WriteString(key, val string) {
	s.writeKey(key)
	s.sb.WriteRune('"')
	s.sb.WriteString(escape(val))
	s.sb.WriteRune('"')
}

func (s *JsonObjectBuilder) WriteInt(key string, val int) {
	s.writeKey(key)
	s.sb.WriteString(strconv.Itoa(val))
}

func (s *JsonObjectBuilder) writeKey(key string) {
	if s.keyCount > 0 {
		s.sb.WriteString(", ")
	}
	s.sb.WriteRune('"')
	s.sb.WriteString(key)
	s.sb.WriteString("\": ")
	s.keyCount++
}

var escapeLookup = [128]string{'\b': "\\b", '\f': "\\f", '\n': "\\n", '\r': "\\r", '\t': "\\t", '"': `\"`, '\\': `\\`}

func escape(s string) string {
	var sb strings.Builder
	hasMapped := false

	for i, r := range s {
		if int(r) < len(escapeLookup) && escapeLookup[r] != "" {
			if !hasMapped {
				sb.Grow(len(s) + 5)
				sb.WriteString(s[:i])
				hasMapped = true
			}
			sb.WriteString(escapeLookup[r])
		} else if hasMapped {
			sb.WriteRune(r)
		}
	}

	if hasMapped {
		return sb.String()
	}
	return s
}

func isNumeric(s string) bool {
	for _, r := range s {
		if (r < '0' || r > '9') && r != '.' {
			return false
		}
	}
	return true
}
