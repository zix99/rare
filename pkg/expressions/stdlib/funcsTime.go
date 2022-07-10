package stdlib

import (
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	. "rare/pkg/expressions" //lint:ignore ST1001 Legacy

	"github.com/araddon/dateparse"
)

// Global to allow changing default time formatting
var DefaultLocalTime = false

const defaultTimeFormat = time.RFC3339

var timeFormats = map[string]string{
	// Standard formats
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
	"MONTHNAME": "January",
	"MONTHA":    "Jan", // a for abbreviated
	"DAY":       "02",
	"YEAR":      "2006",
	"HOUR":      "15",
	"MINUTE":    "04",
	"SECOND":    "05",
	"TIMEZONE":  "MST",
	"NTIMEZONE": "-0700",
	"NTZ":       "-0700",
}

// namedTimeFormatToFormat converts a string to a go-format. If not listed above, assumes the string is the format
func namedTimeFormatToFormat(f string) string {
	if mapped, ok := timeFormats[strings.ToUpper(f)]; ok {
		return mapped
	}
	return f
}

// smartDateParseWrapper wraps different types of date parsing and manipulation into a stage
func smartDateParseWrapper(format string, dateStage KeyBuilderStage, f func(time time.Time) string) KeyBuilderStage {
	switch strings.ToLower(format) {
	case "auto": // Auto will attempt to parse every time
		return KeyBuilderStage(func(context KeyBuilderContext) string {
			strTime := dateStage(context)
			val, err := dateparse.ParseAny(strTime)
			if err != nil {
				return ErrorParsing
			}
			return f(val)
		})

	case "": // Empty format will auto-detect on first successful entry
		var atomicFormat atomic.Value
		atomicFormat.Store("")

		return KeyBuilderStage(func(context KeyBuilderContext) string {
			strTime := dateStage(context)
			if strTime == "" { // This is important for future optimization efforts (so an empty string won't be remembered as a valid format)
				return ErrorParsing
			}

			liveFormat := atomicFormat.Load().(string)
			if liveFormat == "" {
				// This may end up run by a few different threads, but it comes at the benefit
				// of not needing a mutex
				var err error
				liveFormat, err = dateparse.ParseFormat(strTime)
				if err != nil {
					return ErrorParsing
				}
				atomicFormat.Store(liveFormat)
			}

			val, err := time.Parse(liveFormat, strTime)
			if err != nil {
				return ErrorParsing
			}
			return f(val)
		})

	default: // non-empty; Set format will resolve to a go date
		parseFormat := namedTimeFormatToFormat(format)
		return KeyBuilderStage(func(context KeyBuilderContext) string {
			strTime := dateStage(context)
			val, err := time.Parse(parseFormat, strTime)
			if err != nil {
				return ErrorParsing
			}
			return f(val)
		})
	}
}

// Parse time into standard unix epoch time (easier to use)
// By default, will attempt to auto-detect and cache format
// {func <time> [format]}
func kfTimeParse(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) < 1 {
		return stageError(ErrorArgCount)
	}

	// Special key-words for time (eg "now")
	if val, ok := EvalStaticStage(args[0]); ok {
		switch strings.ToLower(val) {
		case "now":
			now := strconv.FormatInt(time.Now().Unix(), 10)
			return func(context KeyBuilderContext) string {
				return now
			}
		}
	}

	// Specific format denoted
	format := EvalStageIndexOrDefault(args, 1, "")

	return smartDateParseWrapper(format, args[0], func(t time.Time) string {
		return strconv.FormatInt(t.Unix(), 10)
	})
}

// {func <unixtime> [format] [*utc*/local]}
func kfTimeFormat(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) < 1 {
		return stageError(ErrorArgCount)
	}
	format := namedTimeFormatToFormat(EvalStageIndexOrDefault(args, 1, defaultTimeFormat))

	utc, utcOk := parseTimeIsUtc(EvalStageIndexOrDefault(args, 2, ""))
	if !utcOk {
		return stageError(ErrorParsing)
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

func timeBucketToFormat(name string) string {
	name = strings.ToLower(name)

	if isPartialString(name, "nanos") {
		return "2006-01-02 15:04:05.999999999"
	} else if isPartialString(name, "seconds") {
		return "2006-01-02 15:04:05"
	} else if isPartialString(name, "minutes") {
		return "2006-01-02 15:04"
	} else if isPartialString(name, "hours") {
		return "2006-01-02 15"
	} else if isPartialString(name, "days") {
		return "2006-01-02"
	} else if isPartialString(name, "months") {
		return "2006-01"
	} else if isPartialString(name, "years") {
		return "2006"
	}
	return ErrorBucket
}

func kfBucketTime(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) < 2 {
		return stageError(ErrorArgCount)
	}

	bucketFormat := timeBucketToFormat(EvalStageOrDefault(args[1], "day"))
	parseFormat := EvalStageIndexOrDefault(args, 2, "")

	return smartDateParseWrapper(parseFormat, args[0], func(t time.Time) string {
		return t.Format(bucketFormat)
	})
}

var attrType = map[string](func(t time.Time) string){
	"WEEKDAY": func(t time.Time) string { return strconv.Itoa(int(t.Weekday())) },
	"WEEK": func(t time.Time) string {
		_, week := t.ISOWeek()
		return strconv.Itoa(week)
	},
	"YEARWEEK": func(t time.Time) string {
		year, week := t.ISOWeek()
		return strconv.Itoa(year) + "-" + strconv.Itoa(week)
	},
	"QUARTER": func(t time.Time) string {
		month := int(t.Month())
		return strconv.Itoa(month/3 + 1)
	},
}

// {func <time> <attr> [utc/local]}
func kfTimeAttr(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) < 2 || len(args) > 3 {
		return stageError(ErrorArgCount)
	}

	attrName, hasAttrName := EvalStaticStage(args[1])
	if !hasAttrName {
		return stageError(ErrorType)
	}
	utc, utcOk := parseTimeIsUtc(EvalStageIndexOrDefault(args, 2, ""))
	if !utcOk {
		return stageError(ErrorParsing)
	}

	attrFunc, hasAttrFunc := attrType[strings.ToUpper(attrName)]
	if !hasAttrFunc {
		return stageError(ErrorBucket)
	}

	return KeyBuilderStage(func(context KeyBuilderContext) string {
		unixTime, err := strconv.ParseInt(args[0](context), 10, 64)
		if err != nil {
			return ErrorType
		}

		t := time.Unix(unixTime, 0)

		if utc {
			t = t.UTC()
		}

		return attrFunc(t)
	})
}

// Pass in "", "local" or "utc"
func parseTimeIsUtc(tzf string) (utc, ok bool) {
	switch strings.ToUpper(tzf) {
	case "":
		return !DefaultLocalTime, true
	case "LOCAL":
		return false, true
	case "UTC":
		return true, true
	default:
		return true, false
	}
}
