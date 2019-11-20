# Expressions

*rare* expressions are handlebars-like in their ability to process data with 
helpers, and format to a new view of the data.  We chose not to do straight
handlebars or templating due to the performance concern that comes with it 
(Processing millions upon millions of templates isn't the use-case html 
processors were meant for)

## Syntax

The syntax for rare-expressions looks like this: `{1} {bucket {2} 100}`.

The basic syntax structure is as follows:

 * Anything not surrounded by `{}` is a literal
 * Expressions are surrounded by `{}`. The entire match will always be `{0}`
 * An integer in an expression denotes a matched value from the regex (or other input) eg. `{2}`
 * When an expression has space(s), the first literal will be the name of a helper function.
   From there, the logic is nested. eg `{coalesce {4} {3} notfound}`
 * Truthiness is the presence of a value.  False is an empty value (or only whitespace)

## Examples

### Parsing an nginx access.log file

Command: `rare histo -m '"(\w{3,4}) ([A-Za-z0-9/.@_-]+).*" (\d{3}) (\d+)' -e "{1} {2} {bytesize {bucket {4} 10000}}" -i "{lt {4} {multi 1024 1024}}" -b access.log`

The above parses the method `{1}`, url `{2}`, status `{3}`, and response size `{4}` in the regex.

It extracts the `<method> <url> <bytesize bucketed to 10k>`. It will ignore `-i` if response size `{4}` is less-than `1024*1024` (1MB).

# Functions

## Coalesce

Syntax: `{coalesce ...}`

Evaluates arguments in-order, chosing the first non-empty result.

## Bucket

Syntax: `{bucket intVal bucketSize}`

Given a value, create equal-sized buckets and place each value in those buckets

## ExpBucket

Syntax: `{expbucket intVal}`

Create exponentially (base-10) increase buckets.

## ByteSize

Syntax: `{bytesize intVal}`

Create a human-readable byte-size format (eg 1024 = 1KB)

## SumI, SubI, MultI, DivI

Syntax: `{sumi ...}`, `{subi ...}`, `{multi ...}`, `{divi ...}`

Evaluates using operator from left to right. Requires at least 2 arguments.

Eg: `{sumi 1 2 3}` will result in `6`

## Equals, NotEquals, Not

Syntax: `{eq a b}`, `{neq a b}`, `{not a}`

Uses truthy-logic to evaluate equality.
eq:  If a == b,  will return "1", otherwise ""
neq: If a != b,  will return "1", otherwise ""
not: If a == "", will return "1", otherwise ""

## LessThan, GreaterThan, LessThanEqual, GreaterThanEqual

Syntax: `{lt a b}`, `{gt a b}`, `{lte a b}`, `{gte a b}`

Uses truthy-logic to compare two integers.

## And, Or

Syntax: `{and ...}`, `{or ...}`

Uses truthy logic and applies `and` or `or` to the values.

and: All arguments need to be truthy
or:  At least one argument needs to be truthy

## Like, Prefix, Suffix

Syntax: `{like val contains}`, `{prefix val startsWith}`, `{suffix val endsWith}`

Truthy check if a value contains a sub-value, starts with, or ends with

## Format

Syntax: `{format "%s" ...}`

Formats a string based on `fmt.Sprintf`

## Humanize Number (Add Commas)

Syntax: `{hf val}`, `{hi val}`

hf: Float
hi: Int

Formats a number based with appropriate placement of commas and decimals

## Tab

Syntax: `{tab a b c ...}`

Concatenates the values of the arguments separated by a table character.

Good for tabulate output separation.

## Json

Syntax: `{json field expression}`

Extract a JSON value based on the expression statement from [gjson](https://github.com/tidwall/gjson)

See: [json](json.md) for more information.

## Time

Syntax: `{time str [format]}` `{timeformat unixtime [format] [utc]}` `{duration dur}`

These three time functions provide you a way to parse and manipulate time.

time: Parse a given time-string into a unix second time (default: RFC3339)
timeformat: Takes a unix time, and formats it (default: RFC3339)
duration: Use a duration expressed in s,m,h and convert it to seconds eg `{duration 24h}`

*Supported Formats:*
ASNIC, UNIX, RUBY, RFC822, RFC822Z, RFC1123, RFC1123Z, RFC3339, RFC3339, RFC3339N, NGINX

*Additional formats for formatting:*
MONTH, DAY, YEAR, HOUR, MINUTE, SECOND, TIMEZONE, NTIMEZONE

*Custom formats:*
You can provide a custom format using go's well-known date. Here's an exercept from their docs:

```
To define your own format, write down what the reference time would look like formatted your way; see the values of constants like ANSIC, StampMicro or Kitchen for examples. The model is to demonstrate what the reference time looks like so that the Format and Parse methods can apply the same transformation to a general time value.

The reference time used in the layouts is the specific time:
Mon Jan 2 15:04:05 MST 2006
```
