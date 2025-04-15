# Expressions

*rare* expressions are handlebars-like in their ability to process data with 
helpers, and format to a new view of the data.  We chose not to do straight
handlebars or go-templating due to the performance concern that comes with it 
(Processing millions upon millions of templates isn't the use-case html 
processors were meant for)

## Syntax

The syntax for rare-expressions looks like this: `{1} {bucket {2} 100}`.

The basic syntax structure is as follows:

 * Anything not surrounded by `{}` is a literal
 * Characters can be escaped with `\`, including `\{` or `\n`
 * Expressions are surrounded by `{}`.
 * An integer in an expression denotes a matched value from the regex (or other input) eg. `{2}`. The entire match will always be `{0}`
 * A string in an expression is a special key or a named regex/dissect group eg. `{src}` or `{group1}`
 * When an expression has space(s), the first literal will be the name of a helper function.
   From there, the logic is nested. eg `{coalesce {4} {3} notfound}`
 * Quotes in an argument create a single argument eg. `{coalesce {4} {3} "not found"}`
 * Truthiness is the presence of a value.  False is an empty value (or only whitespace)

### Special Keys

The following are special Keys:

 * `{src}`  The source name (eg filename). `stdin` when read from stdin
 * `{line}` The line numbers of the current match
 * `{.}`    Returns all matched values with match names as JSON
 * `{#}`    Returns all matched numbered values as JSON
 * `{.#}`   Returned numbered and named matches as JSON
 * `{@}`    All extracted matches in array form

### Testing

You can test and benchmark expressions with the `rare expression` command. For example

```sh
$ rare expression -d 15 -d 20 -k key=30 "The sum is {sumi {0} {1} {key}}"
The sum is 65
```

### Extending

You can add custom functions to *rare* using a [Funcs File](funcsfile.md)

## Examples

### Parsing an nginx access.log file

Command:
```
rare histo \
	-m '"(\w{3,4}) ([A-Za-z0-9/.@_-]+).*" (\d{3}) (\d+)' \
	-e "{1} {2} {bytesize {bucket {4} 10000}}" \
	-i "{lt {4} {multi 1024 1024}}" \
	-b access.log
```

The above parses the method `{1}`, url `{2}`, status `{3}`, and response size `{4}` in the matcher.

It extracts the `<method> <url> <bytesize bucketed to 10k>`. It will ignore `-i` if response size `{4}` is less-than `1024*1024` (1MB).

### Parsing nginx into named groups

Command:
```sh
rare histo \
	-m '"(?P<method>\w{3,4}) (?P<url>[A-Za-z0-9/.@_-]+).*" (?P<status>\d{3}) (?P<size>\d+)' \
	-e "{method} {url} {bytesize {bucket {size} 10000}}" \
	-b access.log
```

In addition to extracting the same number-groups as above, in this case, it will also extract the named-keys
of `{method}`, `{url}`, `{status}`, and `{size}`.

## Functions

**Note on literals:** Some functions take in constant/literals as 1 or more arguments
within the expressions.  These literals will be evaluated during compile-time, rather
than aggregation-time, and will be treated as a literal.  They are denoted below in quotes.

Arguments surrounded by `[]` are optional.

### Aggregation

#### Coalesce

Syntax: `{coalesce ...}`

Evaluates arguments in-order, choosing the first non-empty result.

#### Select Field

Syntax: `{select {0} 1}`

Assuming that `{0}` is a whitespace-separated value, split the values and select the item at index `1`

Eg. `{select "ab cd ef" 1}` will result in `cd`

#### Bucket

Syntax: `{bucket intVal "bucketSize"}`

Given a value, create equal-sized buckets and place each value in those buckets

eg. `{bucket 70 50}` will return `50`

#### BucketRange

Syntax: `{bucketrange intVal "bucketSize"}`

Given a value, create equal-sized buckets and place value into bucket. Outputs range of bucket.

eg. `{bucketrange 70 50}` will return `50 - 99`

#### ExpBucket

Syntax: `{expbucket intVal}`

Create exponentially (base-10) increase buckets.

### Arithmetic

#### Math Formulas

Syntax: `{! "expr"}`

Evaluates a mathematic expression, optionally referencing matches.

Variables in expressions are referenced using brackets rather than
braces. eg. `[0]` instead of `{0}`

Eg: `{! 2+2}`, `{! 2+x}`, `{! 2+[0]}`

See: [Math](math.md) for more information

#### Sumi, Subi, Multi, Divi, Modi

Syntax: `{sumi ...}`, `{subi ...}`, `{multi ...}`, `{divi ...}`, `{modi ...}`

Evaluates integers using operator from left to right. Requires at least 2 arguments.

Eg: `{sumi 1 2 3}` will result in `6`

#### Maxi, Mini

Syntax: `{maxi ...}`, `{mini ...}`

Picks the larger or smallest integer, respectively

Eg: `{mini 3 4 1}` will result in `1`

#### Sumf, Subf, Multf, Divf

Syntax: `{sumf ...}`, `{subf ...}`, `{multf ...}`, `{divf ...}`

Evaluates floating points using operator from left to right. Requires at least 2 arguments.

Eg: `{sumf 1 2 3}` will result in `6`

#### Floor, Ceil, Round

Syntax: `{floor val}`, `{ceil val}`, `{round val [precision=0]}`

Returns the floor, ceil, or rounded format of a floating-point number.

Eg: `{floor 123.765}` will result in `123`

#### Log, Pow, Sqrt

Syntax: `{log10 val}`, `{log2 val}`, `{ln val}`, `{pow val exp}`, `{sqrt val}`

Returns the log (10, 2, or natural), power, or sqrt of a floating-point number.

#### Clamp

Syntax: `{clamp intVal "min" "max"}`

Clamps a given input `intVal` between `min` and `max`.  If falls outside bucket, returns
the word "min" or "max" as appropriate.  If you wish to not see these values, you can
filter with `--ignore`

#### Len

Syntax: `{len string}`

Returns the length of the provided `string`. eg. the string of `hello` returns 5.

### Logic

#### If, Unless

Syntax: `{if val ifTrue ifFalse}`, `{if val ifTrue}`, `{unless val ifFalse}`

If `val` is truthy, then return `ifTrue` else optionally return `ifFalse`

#### Switch

Syntax: `{switch ifTrue val ifTrue val ... [ifFalseVal]}`

In pairs, if a given value is truthy, return the value immediately after. If
there is an odd number of arguments, the last value is used as the "else" result.
Otherwise, empty string is returned.

#### Equals, NotEquals, Not

Syntax: `{eq a b}`, `{neq a b}`, `{not a}`

Uses truthy-logic to evaluate equality.

 * eq:  If a == b,  will return "1", otherwise ""
 * neq: If a != b,  will return "1", otherwise ""
 * not: If a == "", will return "1", otherwise ""

#### LessThan, GreaterThan, LessThanEqual, GreaterThanEqual

Syntax: `{lt a b}`, `{gt a b}`, `{lte a b}`, `{gte a b}`

Uses truthy-logic to compare two integers.

#### And, Or

Syntax: `{and ...}`, `{or ...}`

Uses truthy logic and applies `and` or `or` to the values.

 * and: All arguments need to be truthy
 * or:  At least one argument needs to be truthy

#### Like, Prefix, Suffix

Syntax: `{like val contains}`, `{prefix val startsWith}`, `{suffix val endsWith}`

Truthy check if a value contains a sub-value, starts with, or ends with

#### IsInt, IsNum

Syntax: `{isint val}`, `{isnum val}`

Returns truthy if the val is an integer (isint), or a floating point (isnum)

### Formatting
#### Format

Syntax: `{format "%s" ...}`

Formats a string based on `fmt.Sprintf`: [Go Docs](https://pkg.go.dev/fmt)

#### Substring

Syntax: `{substr {0} pos length}`

Takes the substring of the first argument starting at `pos` for `length`

#### Upper, Lower

Syntax: `{upper val}`, `{lower val}`

Converts a string to all-upper or all-lower case

#### Repeat

Syntax: `{repeat "string" {numtimes}}`

Repeats the "string" the specified number of times

#### Humanize Number (Add Commas)

Syntax: `{hf val}`, `{hi val}`

 * hf: Float
 * hi: Int

Formats a number based with appropriate placement of commas and decimals

#### Percent

Syntax: `{percent val ["precision=1"] [[min=0] max=1]}`

Formats a number as a percentage.  By default, assumes the range is 0-1, therefor
`0.1234` becomes `12.3%`.

Eg.

 * `{percent 0.1234}` will result in `12.3%`
 * `{percent 0.1234 2` will result in `12.34%`
 * `{percent 25 0 100}` will result in `25%`
 * `{percent 100 4 50 150}` will result in `50.0000%`

#### ByteSize, ByteSizeSi

Syntax: `{bytesize intVal [precision=0]}`, `{bytesizesi intVal [precision=0]}`

Create a human-readable byte-size format (eg 1024 = 1KB), or in SI units (1000 = 1KB).
An optional precision allows adding decimals.

#### Downscale

Syntax: `{downscale intVal [precision=0]}`

Formats numbers by thousands (k), Millions (M), Billions (B), or Trillions (T).
eg. `{downscale 10000}` will result in `10k`

### Collecting

#### Tab

Syntax: `{tab a b c ...}`

Concatenates the values of the arguments separated by a table character.

#### CSV

Syntax: `{csv a b c}`

Generate a CSV row given a set of values
#### Arrays / Null Separator

Syntax: `{$ a b c}`

Concatenates a set of arguments with a null separator.  Commonly used
to form arrays that have meaning for a given aggregator.

Specifying multiple expressions is equivalent, eg. `{$ a b}` is the same as `-e a -e b`

### File Loading and Lookup Tables

Load external static content as either raw string, or to be used to lookup
a value given a key.

#### Load

Syntax: `{load "filename"}`

Loads a given filename as text.

To globally disable file loading in expressions for security reasons, specify
`--noload` as global argument.

#### Lookup, HasKey

Syntax: `{lookup key "kv-pairs" ["commentPrefix"]}`, `{haskey key "kv-pairs" ["commentPrefix"]}`

Given a set of kv-pairs (eg. from a loaded file), lookup a key. For `lookup` return a value
and for `haskey` return truthy or falsey.

If a `commentPrefix` is provided, lines in lookup text are ignored if they start with the prefix.

Example kv-pairs text. Keys and values are separated by any whitespace.

```
key1 val1
key2 val2
#comment if '#' set as prefix
key3 val3

#blank lines are ignored
too many values are also ignored
```

### Ranges (Arrays)

Range functions provide the ability to work with arrays in expressions. You
can create an array either manually with the `{@ ...}` function or
by `{@split ...}` a string into an array.

#### Array Definition

Syntax: `{@ ele0 ele1 ele2}` (`{$ ele0 ele1 ele2}` is equivalent)

Creates an array with the provided elements. Use `{@}` for an array of all matches.

#### @len

!!! warning
	This is a linear-time operation. Length of the array
	is not stored and the string needs to be scanned.

Syntax: `{@len <arr>}`

Returns the length of an array.  Empty "" returns 0, a literal will be 1.

#### @in

Syntax: `{@in <val> <arr>}` or `{@in <val> {@ val0 val1 val2 ...}}`

Returns truthy if a given `val` is contained within the array.

#### @split

Syntax: `{@split <arr> ["delim"]}`

Splits a string into an array with the separating `delim`.  If `delim` isn't
specified, `" "` will be used.

#### @join

Syntax: `{@join <arr> ["delim"]}`

Re-joins an array back into a string.  If `delim` is empty, it will be `" "`

#### @map

Syntax: `{@map <arr> <mapfunc>}`

Evaluates `mapfunc` against each element in the array. In `mapfunc`, `{0}`
is the current element.  The function must be surrounded by quotes.

For example, given the array `[1,2,3]`, and the function
`{@map {array} "{multi {0} 2}"}` will output [2,4,6].

#### @reduce

Syntax: `{@reduce <arr> <reducefunc> [initial=""]}`

Evaluates `reducefunc` against each element and a memo. `{0}` is the memo, and
`{1}` is the current value.

For example, given the array `[1,2,3]`, and the function
`{@reduce {array} "{sumi {0} {1}}"}`, it will return `6`.

If `initial` is unset, it will use `arr[0]` as the initial value.

#### @filter

Syntax: `{@filter <arr> <filterfunc>}`

Evaluates `filterfunc` for each element.  If *truthy*, item will be in resulting
array. If false, it will be omitted. `{0}` will be the value examined.

For example, given the array `[1,abc,23,efg]`, and the function
`{@filter {array} "{isnum {0}}"}` will return `[1,23]`.

#### @select

Syntax: `{@select <arr> "index"}`

Selects a single item at an `index` out of `array`.

#### @slice

Syntax: `{@slice <arr> "begin" ["length"]}`

Gets a slice of an array. If `begin` is a negative number, will start from the end.

Examples: (Array `[1,2,3,4]`)

- `{@slice {array} 1}` - [2,3,4]
- `{@slice {array} 1 1}` - [2]
- `{@slice {array} -2}` - [3,4]
- `{@slice {array} -2 1}` - [3]


#### @range

!!! warning
	Since `range` creates an array, large arrays will consume
	a lot of memory. For non-static arrays, it will be created
	each time and could be slow.

Syntax: `{@range [start=0] <stop> [incr=1]}`

Creates an array from start..stop, incrementing by `incr`. Start
defaults to `0` and incr to `1`

Eg:

`{@range 5}` will result in `[0, 1, 2, 3, 4]`

`{@range 1 10 2}` will result in `[1, 3, 5, 7, 9]`


#### @for

!!! warning
	Since `for` creates an array, large arrays will consume
	a lot of memory. For non-static arrays, it will be created
	each time and could be slow.

Syntax: `{@for <start> <whileExpr> <incrExpr>}`

Unlike `@range`, `@for` uses expressions to increment and check when done
as a *truthy* statement.  In the sub-expressions `{0}` is the current value
and `{1}` is the index of the increment.

Eg.

`{@for 0 {lt {0} 5} {sumi {0} 1}}` will result in `[0, 1, 2, 3, 4]`

or something more complex, such as a doubling sequence:

`{@for 1 {lt {1} 5} {sumi {0} {0}}}` will result in `[1, 2, 4, 8, 16]`

### Drawing

#### Colors

Syntax: `{color "color" {string}}`

Available colors: Black, Red, Green, Yellow, Blue, Magenta, Cyan, White

Note: If colors are disabled, no color will be shown.

Colorizes the 2nd argument.

#### Bars

Syntax: `{bar {val} "maxVal" "length" ["scale"]}`

Note: If unicode is disabled, will use pipe character

Draws a "bar" with the length `(val / maxVal) * length`

Scale can be `linear`, `log10`, or `log2` (Default: `lienar`)

### Paths

Syntax: `{basename a/b/c}`, `{dirname a/b/c}`, `{extname a/b/c.jpg}`

Selects the base, directory, or extension of a path.

 * `basename a/b/c` = c
 * `dirname  a/b/c` = a/b
 * `extname a/b/c.jpg` = .jpg 

### Json

Syntax: `{json field expression}` or `{json expression}`

Extract a JSON value based on the expression statement from [gjson](https://github.com/tidwall/gjson)

When only 1 argument is present, it will assume the JSON is in `{0}` (Full match)

See: [json](json.md) for more information.

### Time

#### Time Parsing

Syntax: `{time str "[format:cache]" "[tz:utc]"}`

Parse a given time-string into a unix second time (default format: `cache`)

##### Format Auto-Detection

If the format argument is omitted or set to "auto", it will attempt to resolve the format of the time.

If the format is unable to be resolved, it must be specified manually with a format below, or a custom format.

If omitted or "cache": The first seen date will determine the format for all dates going forward (faster)

If "auto": The date format will be auto-detected with each parse. This can be used if the date could be in different formats (slower)

##### Timezones

The following values are accepted for a `tz` (timezone): `utc`, `local`, or a valid *IANA Time Zone*

By default, all datetimes are processed as UTC, unless explicit in the datetime itself, or overridden via a parameter.


#### Time Values

Syntax: `{time now}`

These are special values to output:

 * `now` - return the current unix timestamp cached at the start of *rare*
 * `live` - return the current unix timestamp the moment it's executed
 * `delta` - return the number of seconds that *rare* has executed


#### Time Format

Syntax: `{timeformat unixtime "[format:RFC3339]" "[tz:utc]"}`

Takes a unix time, and formats it (default: `RFC3339`)

To reformat a time, you need to parse it first, eg: `{timeformat {time {0}} RFC3339}`

**Supported Formats:**
ANSIC, UNIX, RUBY, RFC822, RFC822Z, RFC1123, RFC1123Z, RFC3339, RFC3339, RFC3339N, NGINX

**Additional formats for formatting:**
MONTH, MONTHNAME, MNTH, DAY, WEEKDAY, WDAY, YEAR, HOUR, MINUTE, SECOND, TIMEZONE, NTIMEZONE

**Custom formats:**
You can provide a custom format using go's well-known date. Here's an exercept from the docs:

To define your own format, write down what the reference time would look like formatted your way; see the values of constants
like ANSIC, StampMicro or Kitchen for examples. The model is to demonstrate what the reference time looks like so that the Format
and Parse methods can apply the same transformation to a general time value.

The reference time used in the layouts is the specific time: `Mon Jan 2 15:04:05 MST 2006`

#### Time Attribute

Syntax: `{timeattr unixtime attr [tz:utc]"}`

Extracts an attribute about a given datetime

*Supports*: `weekday`, `week`, `yearweek`, `quarter`

#### Durations

Syntax: `{duration dur}`

Use a duration expressed in s,m,h and convert it to seconds eg `{duration 24h}`

From there, you can do arithmetic on time, for instance: `{sumi {time now} {duration 1h}}`

##### Format Duration

Syntax: `{durationformat secs}`

Formats a duration (in seconds) to a human-readable time, (eg. 4h0m0s)

#### Time Bucket

Syntax: `{buckettime str bucket "[format]" "[tz:utc]"}`

Truncate the time to a given bucket (*n*ano, *s*econd, *m*inute, *h*our, *d*ay, *mo*nth, *y*ear)


## Errors

The following error strings may be returned while compiling or evaluating your expression

| Name      | Code            | Description                                                               |
|-----------|-----------------|---------------------------------------------------------------------------|
| Type      | `<BAD-TYPE>`    | Error parsing the principle value of the input because of unexpected type |
| Parsing   | `<PARSE-ERROR>` | Error parsing the principle value of the input (non-numeric)              |
| Arg Count | `<ARGN>`        | Function to not support a variation with the given argument count         |
| Const     | `<CONST>`       | Expected constant value                                                   |
| Enum      | `<ENUM>`        | A given value is not contained within a set                               |
| Arg Name  | `<NAME>`        | A variable accessed by a given name does not exist                        |
| Empty     | `<EMPTY>`       | A value was expected, but was empty                                       |
| File      | `<FILE>`        | Unable to read file                                                       |
| Value     | `<VALUE>`       | Value is out of range or invalid (eg. range incrementer is 0)             |
