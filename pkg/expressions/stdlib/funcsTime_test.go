package stdlib

import (
	"rare/pkg/expressions"
	"strconv"
	"testing"
	"time"

	"github.com/araddon/dateparse"
	"github.com/stretchr/testify/assert"
)

func TestTimeExpression(t *testing.T) {
	testExpression(t,
		mockContext("14/Apr/2016:19:12:25 +0200"),
		"{time {0} NGINX}",
		"1460653945")
}

func TestFormatExpression(t *testing.T) {
	testExpression(t,
		mockContext("14/Apr/2016:19:12:25 +0200"),
		"{timeformat {time {0} NGINX} RFC3339 utc}",
		"2016-04-14T17:12:25Z")
}

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

func TestTimeFormat(t *testing.T) {
	testExpression(t,
		mockContext("14/Apr/2016:19:12:25 +0200"),
		"{buckettime {0} d nginx}",
		"2016-04-14")
}

func TestTimeFormatDetection(t *testing.T) {
	testExpression(t,
		mockContext("14/Apr/2016:19:12:25 +0200"),
		"{buckettime {0} d}",
		"2016-04-14")
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

func BenchmarkTimeParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Parse(time.RFC3339, "2012-05-02T15:04:05Z07:00")
	}
}

func BenchmarkDateParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dateparse.ParseAny("2012-05-02T15:04:05Z07:00")
	}
}
