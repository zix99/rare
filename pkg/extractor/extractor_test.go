package extractor

import (
	"rare/pkg/matchers"
	"rare/pkg/matchers/fastregex"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testData = `abc 123
def 245
qqq 123
xxx`

func TestBasicExtractor(t *testing.T) {
	input := convertReaderToBatches("test", strings.NewReader(testData), 1)
	ex, err := New(input, &Config{
		Matcher: matchers.ToFactory(fastregex.MustCompile(`(\d+)`)),
		Extract: "val:{1}",
		Workers: 1,
	})
	assert.NoError(t, err)

	vals := unbatchMatches(ex.ReadChan())
	assert.Equal(t, "abc 123", vals[0].Line)
	assert.Equal(t, 4, len(vals[0].Indices))
	assert.Equal(t, "val:123", vals[0].Extracted)
	assert.Equal(t, uint64(1), vals[0].LineNumber)
	assert.Equal(t, "test", vals[0].Source)

	assert.Equal(t, 3, len(vals))
	assert.Equal(t, uint64(2), vals[1].LineNumber)

	assert.Equal(t, uint64(0), ex.IgnoredLines())
	assert.Equal(t, uint64(3), ex.MatchedLines())
	assert.Equal(t, uint64(4), ex.ReadLines())
}

func TestSourceAndLine(t *testing.T) {
	input := convertReaderToBatches("test", strings.NewReader(testData), 1)
	ex, err := New(input, &Config{
		Matcher: matchers.ToFactory(fastregex.MustCompile(`(\d+)`)),
		Extract: "{src} {line} val:{1} {bad} {@}",
		Workers: 1,
	})
	assert.NoError(t, err)

	vals := unbatchMatches(ex.ReadChan())
	assert.Equal(t, "test 1 val:123 <NAME> 123", vals[0].Extracted)
	assert.Equal(t, uint64(1), vals[0].LineNumber)

	assert.Equal(t, "test 2 val:245 <NAME> 245", vals[1].Extracted)
	assert.Equal(t, "test 3 val:123 <NAME> 123", vals[2].Extracted)
}

func TestIgnoreLines(t *testing.T) {
	input := convertReaderToBatches("test", strings.NewReader(testData), 1)
	ignore, _ := NewIgnoreExpressions(`{eq {1} "123"}`)
	ex, err := New(input, &Config{
		Matcher: matchers.ToFactory(fastregex.MustCompile(`(\d+)`)),
		Extract: "{src} {line} val:{1} {bad}{500}",
		Workers: 1,
		Ignore:  ignore,
	})
	assert.NoError(t, err)

	vals := unbatchMatches(ex.ReadChan())

	assert.Len(t, vals, 1)
}

func TestNamedGroup(t *testing.T) {
	input := convertReaderToBatches("test", strings.NewReader(testData), 1)
	ex, err := New(input, &Config{
		Matcher: matchers.ToFactory(fastregex.MustCompile(`(?P<num>\d+)`)),
		Extract: "val:{1}:{num}",
		Workers: 1,
	})
	assert.NoError(t, err)

	vals := unbatchMatches(ex.ReadChan())
	assert.Equal(t, "abc 123", vals[0].Line)
	assert.Equal(t, 4, len(vals[0].Indices))
	assert.Equal(t, "val:123:123", vals[0].Extracted)
}

func TestJSONOutput(t *testing.T) {
	input := convertReaderToBatches("test", strings.NewReader(testData), 1)
	ex, err := New(input, &Config{
		Matcher: matchers.ToFactory(fastregex.MustCompile(`(?P<num>\d+)`)),
		Extract: "{.} {#} {.#} {#.}",
		Workers: 1,
	})

	assert.NoError(t, err)
	vals := unbatchMatches(ex.ReadChan())
	assert.Equal(t, `{"num": 123} {"0": 123, "1": 123} {"num": 123, "0": 123, "1": 123} {"num": 123, "0": 123, "1": 123}`, vals[0].Extracted)
}

func TestGH10SliceBoundsPanic(t *testing.T) {
	input := convertReaderToBatches("", strings.NewReader("this is an [ERROR] message"), 1)
	ex, err := New(input, &Config{
		Matcher: matchers.ToFactory(fastregex.MustCompile(`\[(INFO)|(ERROR)|(WARNING)|(CRITICAL)\]`)),
		Extract: "val:{2} val:{3}",
		Workers: 1,
	})
	assert.NoError(t, err)

	vals := unbatchMatches(ex.ReadChan())
	assert.Equal(t, "val:ERROR val:", vals[0].Extracted)
	assert.Equal(t, []int{12, 17, -1, -1, 12, 17, -1, -1, -1, -1}, vals[0].Indices)
}
