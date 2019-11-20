package expressions

import (
	"strconv"
	"testing"
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
		"{timeformat {time {0} NGINX}}",
		"2016-04-14T13:12:25-04:00")
}

func TestAddDurationDay(t *testing.T) {
	testExpression(t,
		mockContext("14/Apr/2016:19:12:25 +0200"),
		"{timeformat {sumi {time {0} NGINX} {duration 24h}}}",
		"2016-04-15T13:12:25-04:00")
}

func TestDuration(t *testing.T) {
	testExpression(t,
		mockContext(),
		"{duration 24h}",
		strconv.Itoa(60*60*24))
}
