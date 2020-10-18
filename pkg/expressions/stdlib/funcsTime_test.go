package stdlib

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
