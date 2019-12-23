package expressions

import (
	"strconv"
	"strings"
	"time"
)

const defaultTimeFormat = time.RFC3339

var timeFormats = map[string]string{
	"ASNIC":    time.ANSIC,
	"UNIX":     time.UnixDate,
	"RUBY":     time.RubyDate,
	"RFC822":   time.RFC822,
	"RFC822Z":  time.RFC822Z,
	"RFC1123":  time.RFC1123,
	"RFC1123Z": time.RFC1123Z,
	"RFC3339":  time.RFC3339,
	"RFC3339N": time.RFC3339Nano,
	// Custom formats,
	"NGINX": "_2/Jan/2006:15:04:05 -0700",
	// Parts,
	"MONTH":     "01",
	"DAY":       "_2",
	"YEAR":      "2006",
	"HOUR":      "15",
	"MINUTE":    "04",
	"SECOND":    "05",
	"TIMEZONE":  "MST",
	"NTIMEZONE": "-0700",
	"NTZ":       "-0700",
}

func namedTimeFormatToFormat(f string) string {
	if mapped, ok := timeFormats[strings.ToUpper(f)]; ok {
		return mapped
	}
	return f
}

// Parse time into standard unix time (easier to use)
func kfTimeParse(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) < 1 {
		return stageError(ErrorArgCount)
	}
	format := defaultTimeFormat
	if len(args) >= 2 {
		format = namedTimeFormatToFormat(args[1](nil))
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		strTime := args[0](context)
		val, err := time.Parse(format, strTime)
		if err != nil {
			return ErrorParsing
		}
		return strconv.FormatInt(val.Unix(), 10)
	})
}

func kfTimeFormat(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) < 1 {
		return stageError(ErrorArgCount)
	}
	format := defaultTimeFormat
	if len(args) >= 2 {
		format = namedTimeFormatToFormat(args[1](nil))
	}
	utc := false
	if len(args) >= 3 {
		utc = Truthy(args[2](nil))
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		strUnixTime := args[0](context)
		unixTime, err := strconv.ParseInt(strUnixTime, 10, 64)
		if err != nil {
			return ErrorType
		}
		t := time.Unix(unixTime, 0)
		if utc {
			t = t.UTC()
		}
		return t.Format(format)
	})
}

func kfDuration(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 1 {
		return stageError(ErrorArgCount)
	}

	return KeyBuilderStage(func(context KeyBuilderContext) string {
		strDuration := args[0](context)

		duration, err := time.ParseDuration(strDuration)
		if err != nil {
			return ErrorType
		}

		return strconv.FormatInt(int64(duration.Seconds()), 10)
	})
}
