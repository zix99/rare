# Extractor (Matcher)

The main component of *rare* is the extractor (or matcher).  There are
three fundamental concepts around the parser:

 * Each line of an input (separated by `\n`) is matched to a matcher
 * A matcher is used to parse a line into a match (and optionally, groups)
 * An expression (see: [expression](expressions.md)) is used to format an
   output from a matched groups
 * Optionally, one or more ignore expressions can be applied to silent matches
   that satisfy a truthy-comparison

## Matcher Types

If no matcher is specified, by default, the entire line is always matched
and passed-through to the expression-stage.

More than one matcher can **not** be specified at the same time.

### Regex

A regex expression is specified with `--match` or `-m`, and follows common
[regex syntax](regexp.md).

When matching a regex, groups and keys are extracted both index and
by-name if specified.

Set ignore-case with `-I` or `--ignore-case`.

**Example:**

```bash
rare filter -m '"(\w{3,4}) ([A-Za-z0-9/.@_-]+)' access.log
```

### Dissect

A dissect expression is specified with `--disect` or `-d`, and follows
[dissect syntax](dissect.md).

Like regex, groups are extracted by both index and name. Dissect can
be significantly faster than regex.

Set ignore-case with `-I` or `--ignore-case`.

**Example:**

```bash
rare filter -d 'HTTP/1.1" %{code} %{size} ' -e '{code}' access.log
```

## Ignore

You can provide one or more [expressions](expressions.md) via `--ignore` (`-i`). If
the statement evaluates to truthy (non-empty), the matched line will be ignored.

**Example:**

To ignore all non-`200` http codes

```bash
rare filter -d 'HTTP/1.1" %{code} %{size} ' -i '{neq {code} 200}' -e '{code} {size}' access.log
```

## Examples

### Decomposing a Matcher

The most primitive way use rare is to filter lines in an input.  We'll
be using an example nginx log for our example.

Nginx log line looks like this:

```log
10.20.30.40 - - [14/Apr/2016:18:13:29 +0200] "GET / HTTP/1.1" 206 101 "-" "curl/7.43.0"
```

So, to parse this we may want to match the method and path with a regex:

```bash
rare filter -m '"(\w{3,4}) ([A-Za-z0-9/.@_-]+)' access.log
```

This will extract the method and url and output the entire line to the screen (if matched).

If you want it to only output the matched portion, you can add `-e "{0}"`

Lastly, lets say we want to ignore all paths that equal "/", we could do that by adding
an ignore pattern: `-i {eq {1} /}`

### Histograms

Histograms are like filters, but rather than outputting every match, it will
create an aggregated count based on the extracted expression.

So, with the same example as above, if we extract the method and url, we will
get something that will count based on keys.

```bash
rare histogram -m '"(\w{3,4}) ([A-Za-z0-9/.@_-]+)' -e '{1} {2}' -b access.log
```

## See Also

* [Regular Expressions](regexp.md)
* [Examples](examples.md)