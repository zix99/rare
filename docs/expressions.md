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
