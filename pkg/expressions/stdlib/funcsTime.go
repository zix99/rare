package stdlib

import (
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	. "rare/pkg/expressions" //lint:ignore ST1001 Legacy

	"github.com/araddon/dateparse"
)

const defaultTimeFormat = time.RFC3339

var timeFormats = map[string]string{
	"": defaultTimeFormat,
	// Standard formats
	"ANSIC":    time.ANSIC,
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
	"MNTH":      "Jan",
	"DAY":       "02",
	"YEAR":      "2006",
	"HOUR":      "15",
	"MINUTE":    "04",
	"SECOND":    "05",
	"TIMEZONE":  "MST",
	"NTIMEZONE": "-0700",
	"NTZ":       "-0700",
	"WEEKDAY":   "Monday",
	"WDAY":      "Mon",
}

// namedTimeFormatToFormat converts a string to a go-format. If not listed above, assumes the string is the format
func namedTimeFormatToFormat(f string) string {
	if mapped, ok := timeFormats[strings.ToUpper(f)]; ok {
		return mapped
	}
	return f
}

// smartDateParseWrapper wraps different types of date parsing and manipulation into a stage
func smartDateParseWrapper(format string, tz *time.Location, dateStage KeyBuilderStage, f func(time time.Time) string) (KeyBuilderStage, error) {
	switch strings.ToLower(format) {
	case "auto": // Auto will attempt to parse every time
		return KeyBuilderStage(func(context KeyBuilderContext) string {
			strTime := dateStage(context)
			val, err := dateparse.ParseIn(strTime, tz)
			if err != nil {
				return ErrorParsing
			}
			return f(val)
		}), nil

	case "", "cache": // Empty format will auto-detect on first successful entry
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

			val, err := time.ParseInLocation(liveFormat, strTime, tz)
			if err != nil {
				return ErrorParsing
			}
			return f(val)
		}), nil

	default: // non-empty; Set format will resolve to a go date
		parseFormat := namedTimeFormatToFormat(format)
		return KeyBuilderStage(func(context KeyBuilderContext) string {
			strTime := dateStage(context)
			val, err := time.ParseInLocation(parseFormat, strTime, tz)
			if err != nil {
				return ErrorParsing
			}
			return f(val)
		}), nil
	}
}

// Parse time into standard unix epoch time (easier to use)
// By default, will attempt to auto-detect and cache format
// {func <time> [format:cache] [tz:utc]}
func kfTimeParse(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if !isArgCountBetween(args, 1, 3) {
		return stageErrArgRange(args, "1-3")
	}

	// Special key-words for time (eg "now")
	if val, ok := EvalStaticStage(args[0]); ok {
		switch strings.ToLower(val) {
		case "now":
			now := strconv.FormatInt(time.Now().Unix(), 10)
			return func(context KeyBuilderContext) string {
				return now
			}, nil
		case "live":
			return func(context KeyBuilderContext) string {
				context.GetMatch(-1) // HACK: Touch the context so it doesn't get optimized out
				return strconv.FormatInt(time.Now().Unix(), 10)
			}, nil
		case "delta":
			start := time.Now().Unix()
			return func(context KeyBuilderContext) string {
				context.GetMatch(-1) // HACK: Touch the context so it doesn't get optimized out
				return strconv.FormatInt(time.Now().Unix()-start, 10)
			}, nil
		}
	}

	// Specific format denoted
	format := EvalStageIndexOrDefault(args, 1, "")
	tz, tzOk := parseTimezoneLocation(EvalStageIndexOrDefault(args, 2, ""))
	if !tzOk {
		return stageArgError(ErrParsing, 2)
	}

	return smartDateParseWrapper(format, tz, args[0], func(t time.Time) string {
		return strconv.FormatInt(t.Unix(), 10)
	})
}

// {func <unixtime> [format:RFC3339] [tz:utc]}
func kfTimeFormat(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if !isArgCountBetween(args, 1, 3) {
		return stageErrArgRange(args, "1-3")
	}
	format := namedTimeFormatToFormat(EvalStageIndexOrDefault(args, 1, defaultTimeFormat))

	tz, tzOk := parseTimezoneLocation(EvalStageIndexOrDefault(args, 2, ""))
	if !tzOk {
		return stageArgError(ErrParsing, 2)
	}

	return KeyBuilderStage(func(context KeyBuilderContext) string {
		strUnixTime := args[0](context)
		unixTime, err := strconv.ParseInt(strUnixTime, 10, 64)
		if err != nil {
			return ErrorNum
		}

		t := time.Unix(unixTime, 0).In(tz)

		return t.Format(format)
	}), nil
}

// {func <duration_string>}
func kfDuration(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 1 {
		return stageErrArgCount(args, 1)
	}

	return KeyBuilderStage(func(context KeyBuilderContext) string {
		strDuration := args[0](context)

		duration, err := time.ParseDuration(strDuration)
		if err != nil {
			return ErrorParsing
		}

		return strconv.FormatInt(int64(duration.Seconds()), 10)
	}), nil
}

// {func <secs>} format seconds to duration
func kfDurationFormat(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 1 {
		return stageErrArgCount(args, 1)
	}

	return KeyBuilderStage(func(context KeyBuilderContext) string {
		secs, err := strconv.ParseInt(args[0](context), 10, 64)
		if err != nil {
			return ErrorNum
		}

		return (time.Duration(secs) * time.Second).String()
	}), nil
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
	return ""
}

// {func <time> <bucket> [format:auto] [tz:utc]}
func kfBucketTime(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if !isArgCountBetween(args, 2, 4) {
		return stageErrArgRange(args, "2-4")
	}

	bucketName, bucketNameOk := EvalStaticStage(args[1])
	if !bucketNameOk {
		return stageArgError(ErrConst, 1)
	}

	bucketFormat := timeBucketToFormat(bucketName)
	if bucketFormat == "" {
		return stageArgError(ErrEnum, 1)
	}

	parseFormat := EvalStageIndexOrDefault(args, 2, "")
	tz, tzOk := parseTimezoneLocation(EvalStageIndexOrDefault(args, 3, ""))
	if !tzOk {
		return stageArgError(ErrParsing, 3)
	}

	return smartDateParseWrapper(parseFormat, tz, args[0], func(t time.Time) string {
		return t.Format(bucketFormat)
	})
}

// valid time attributes
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

// {func <time> <attr> [tz:utc]}
func kfTimeAttr(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) < 2 || len(args) > 3 {
		return stageErrArgRange(args, "2-3")
	}

	attrName, hasAttrName := EvalStaticStage(args[1])
	if !hasAttrName {
		return stageArgError(ErrConst, 1)
	}
	tz, tzOk := parseTimezoneLocation(EvalStageIndexOrDefault(args, 2, ""))
	if !tzOk {
		return stageArgError(ErrParsing, 2)
	}

	attrFunc, hasAttrFunc := attrType[strings.ToUpper(attrName)]
	if !hasAttrFunc {
		return stageArgError(ErrEnum, 1)
	}

	return KeyBuilderStage(func(context KeyBuilderContext) string {
		unixTime, err := strconv.ParseInt(args[0](context), 10, 64)
		if err != nil {
			return ErrorNum
		}

		t := time.Unix(unixTime, 0).In(tz)

		return attrFunc(t)
	}), nil
}

// Pass in "", "local", "utc" or a valid unix timezone
func parseTimezoneLocation(tzf string) (loc *time.Location, ok bool) {
	switch strings.ToUpper(tzf) {
	case "", "UTC":
		return time.UTC, true
	case "LOCAL":
		return time.Local, true
	default:
		if tz, err := time.LoadLocation(tzf); err == nil {
			return tz, true
		}
		return time.UTC, false
	}
}
