package expressions

import (
	"strconv"
	"strings"
	"time"
)

const defaultTimeFormat = time.RFC3339

func namedTimeFormatToFormat(f string) string {
	switch strings.ToUpper(f) {
	// Standard formats
	case "ASNIC":
		return time.ANSIC
	case "UNIX":
		return time.UnixDate
	case "RUBY":
		return time.RubyDate
	case "RFC822":
		return time.RFC822
	case "RFC822Z":
		return time.RFC822Z
	case "RFC1123":
		return time.RFC1123
	case "RFC1123Z":
		return time.RFC1123Z
	case "RFC3339":
		return time.RFC3339
	case "RFC3339N":
		return time.RFC3339Nano
	// Custom formats
	case "NGINX":
		return "_2/Jan/2006:15:04:05 -0700"
	// Parts
	case "MONTH":
		return "01"
	case "DAY":
		return "_2"
	case "YEAR":
		return "2006"
	case "HOUR":
		return "15"
	case "MINUTE":
		return "04"
	case "SECOND":
		return "05"
	case "TIMEZONE":
		return "MST"
	case "NTIMEZONE":
	case "NTZ":
		return "-0700"
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
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		strUnixTime := args[0](context)
		unixTime, err := strconv.ParseInt(strUnixTime, 10, 64)
		if err != nil {
			return ErrorType
		}
		t := time.Unix(unixTime, 0)
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
