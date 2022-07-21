package stdlib

import (
	"rare/pkg/expressions"
	"strconv"
	"testing"
	"time"

	"github.com/araddon/dateparse"
	"github.com/stretchr/testify/assert"
)

// Basic / Time parsing

func TestTimeExpression(t *testing.T) {
	testExpression(t,
		mockContext("14/Apr/2016:19:12:25 +0200"),
		"{time {0} NGINX}",
		"1460653945")
}

func TestFormatExpression(t *testing.T) {
	// Defined type
	testExpression(t,
		mockContext("14/Apr/2016:19:12:25 +0200"),
		"{timeformat {time {0} NGINX} RFC3339 utc}",
		"2016-04-14T17:12:25Z")
	// Explicit
	testExpression(t,
		mockContext("14/Apr/2016:19:12:25 +0200"),
		`{timeformat {time {0} "_2/Jan/2006:15:04:05 -0700"} RFC3339 utc}`,
		"2016-04-14T17:12:25Z")
	// Default/empty-string
	testExpression(t,
		mockContext("14/Apr/2016:19:12:25 +0200"),
		`{timeformat {time {0}} "" utc}`,
		"2016-04-14T17:12:25Z")
}

func TestTimeExpressionDetection(t *testing.T) {
	testExpression(t,
		mockContext("14/Apr/2016:19:12:25 +0200"),
		"{time {0}}",
		"1460653945")
}

func TestTimeNow(t *testing.T) {
	kb, err := NewStdKeyBuilder().Compile("{time now}")
	assert.NoError(t, err)

	val := kb.BuildKey(mockContext())
	assert.NotEmpty(t, val)

	ival, err := strconv.ParseInt(val, 10, 64)
	assert.NoError(t, err)
	assert.NotZero(t, ival)

}

func TestTimeExpressionDetectionFailure(t *testing.T) {
	testExpression(t,
		mockContext("oauef888"),
		"{time {0}}",
		"<PARSE-ERROR>")
}

func TestTimeExpressionDetectionAuto(t *testing.T) {
	testExpression(t,
		mockContext("14/Apr/2016:19:12:25 +0200"),
		"{time {0} auto}",
		"1460653945")
}

// Duration

func TestAddDurationDay(t *testing.T) {
	testExpression(t,
		mockContext("14/Apr/2016:19:12:25 +0200"),
		"{timeformat {sumi {time {0} NGINX} {duration 24h}} RFC822 utc}",
		"15 Apr 16 17:12 UTC")
}

func TestDuration(t *testing.T) {
	testExpression(t,
		mockContext(),
		"{duration 24h}",
		strconv.Itoa(60*60*24))
}

func TestDurationFormat(t *testing.T) {
	testExpression(t,
		mockContext("14400"),
		"{durationformat {0}}",
		"4h0m0s")
	testExpression(t,
		mockContext("14400"),
		"{durationformat {0} b}",
		"<<ARGN>>")
}

// Bucketing

func TestTimeBucketFormat(t *testing.T) {
	testExpression(t,
		mockContext("14/Apr/2016:19:12:25 +0200"),
		"{buckettime {0} d nginx}",
		"2016-04-14")
}

func TestTimeBucketFormatDetection(t *testing.T) {
	testExpression(t,
		mockContext("14/Apr/2016:19:12:25 +0200"),
		"{buckettime {0} d}",
		"2016-04-14")
}

func TestTimeBucketUtc(t *testing.T) {
	testExpression(t,
		mockContext("14/Apr/2016 23:00:00"),
		"{buckettime {0} d}",
		"2016-04-14")
}

// Time attributes

func TestTimeAttr(t *testing.T) {
	testExpression(t,
		mockContext("14/Apr/2016 01:00:00"),
		"{timeattr {time {0}} weekday}",
		"4")
	testExpression(t,
		mockContext("14/Apr/2016 01:00:00"),
		"{timeattr {time {0}} week}",
		"15")
	testExpression(t,
		mockContext("14/Apr/2016 01:00:00"),
		"{timeattr {time {0}} Yearweek}",
		"2016-15")
	testExpression(t,
		mockContext("14/Apr/2016 01:00:00"),
		"{timeattr {time {0}} quarter}",
		"2")
}

func TestTimeAttrToLocal(t *testing.T) {
	kb, err := NewStdKeyBuilder().Compile("{timeattr {time {0}} weekday local}")
	assert.NoError(t, err)
	ret := kb.BuildKey(mockContext("14/Apr/2016 01:00:00"))
	assert.NotEmpty(t, ret)
}

func TestTimeAttrToBadTZ(t *testing.T) {
	testExpression(t,
		mockContext("14/Apr/2016 01:00:00"),
		"{timeattr {time {0}} weekday asdf}",
		"<<PARSE-ERROR>>")
}

func TestTimeAttrArgError(t *testing.T) {
	testExpression(t,
		mockContext("14/Apr/2016 01:00:00"),
		"{timeattr {time {0}}}",
		"<<ARGN>>")
}

func TestTimeAttrArgErrorExtra(t *testing.T) {
	testExpression(t,
		mockContext("14/Apr/2016 01:00:00"),
		"{timeattr {time {0}} a b c}",
		"<<ARGN>>")
}

// Utilities
func TestLoadingTimezone(t *testing.T) {
	tz, ok := parseTimezoneLocation("utc")
	assert.Equal(t, tz, time.UTC)
	assert.True(t, ok)

	tz, ok = parseTimezoneLocation("Local")
	assert.Equal(t, tz, time.Local)
	assert.True(t, ok)

	tz, ok = parseTimezoneLocation("America/New_York")
	assert.NotNil(t, tz)
	assert.True(t, ok)

	tz, ok = parseTimezoneLocation("not a real timezone")
	assert.Equal(t, tz, time.UTC)
	assert.False(t, ok)
}

// BenchmarkTimeParseExpression-4   	  537970	      2133 ns/op	     536 B/op	       9 allocs/op
func BenchmarkTimeParseExpression(b *testing.B) {
	stage := kfTimeParse([]expressions.KeyBuilderStage{
		func(kbc expressions.KeyBuilderContext) string {
			return kbc.GetMatch(0)
		},
		stageLiteral("auto"),
	})
	for i := 0; i < b.N; i++ {
		stage(&expressions.KeyBuilderContextArray{
			Elements: []string{
				"14/Apr/2016:19:12:25 +0200",
			},
		})
	}
}

// BenchmarkTimeParse-4   	 1686390	       654.7 ns/op	     120 B/op	       4 allocs/op
func BenchmarkTimeParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Parse(time.RFC3339, "2012-05-02T15:04:05Z07:00")
	}
}

// BenchmarkDateParse-4   	  757498	      1559 ns/op	     440 B/op	       7 allocs/op
func BenchmarkDateParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dateparse.ParseAny("2012-05-02T15:04:05Z07:00")
	}
}
