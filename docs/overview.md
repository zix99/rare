# rare

Rare is a fast, realtime regex-extraction, and aggregation into common formats
such as histograms, numerical summaries, tables, and more!

Rare is composed of three parts in the pipeline:

1. Extraction
2. Expression Building
3. Aggregation

## Extraction

Extraction is denoted with `-m` and is the process of reading a line in a file or set
of files and parsing it with a regular expression into the match-groups denoted
by the regex.

If the regex doesn't match, the line is discarded (a non-match)

These match groups are then fed into the next stage, the expression.

Read more at:

* rare docs parsing

## Expressions

Expressions `-e` take the match groups, and other pieces of information, and build
a string-based key.  The match groups can be operated on by helpers to build
the string-key (eg arithmatic, json parsing, simple logic).

The result of this key will act as the key for the aggregation stage.

Optionally ignore expressions can be supplied with `-i` which will
cause the match to be ignored if the expression evaluates to be truthy.

Read more at:

* rare docs expressions

## Aggregation

The last stage is the aggregation and realtime presentation.  This takes the
key built in the expression stage and uses it to aggregate on. Different aggregators
aggregate in different ways.  For example, the `histogram` will count instances of the key,
the `table` will count it in 2D, and the `analyze` will treat the match as a number.

Run `rare --help` for more details on the aggregators.
