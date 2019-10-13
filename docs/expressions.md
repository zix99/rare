# Expressions

*rare* expressions are handlebars-like in their ability to process data with helpers, and
format to a new view of the data.  We chose not to do straight handlebars or templating due
to the performance concern that comes with it (Processing millions upon millions of templates
isn't the use-case html processors were meant for)

## Syntax

The syntax for rare-expressions looks like this: `{1} {bucket {2} 100}`.

The basic syntax structure is as follows:

 * Anything not surrounded by `{}` is a literal
 * Expressions are surrounded by `{}`. The entire match will always be `{0}`
 * An integer in an expression denotes a matched value from the regex (or other input) eg. `{2}`
 * When an expression has space(s), the first literal will be the name of a helper function. From there, the logic is nested. eg `{coalesce {4} {3} notfound}`
 * Truthiness is the presence of a value.  False is an empty value (or only whitespace)

## Examples

### Parsing an nginx access.log file

`rare histo -m '"(\w{3,4}) ([A-Za-z0-9/.@_-]+).*" (\d{3}) (\d+)' -e "{1} {2} {bytesize {bucket {4} 10000}}" -i "{lt {4} {multi 1024 1024}}" -b access.log`

The above parses the method `{1}`, url `{2}`, status `{3}`, and response size `{4}` in the regex.

It extracts the `<method> <url> <bytesize bucketed to 10k>`. It will ignore `-i` if response size `{4}` is less-than `1024*1024` (1MB).

# Functions

## Coalesce

## Bucket
