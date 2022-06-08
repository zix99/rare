# rare

Rare is a fast, realtime regex-extraction, and aggregation into common formats
such as histograms, numerical summaries, tables, and more!

Rare is composed of four parts in the pipeline:

1. Batching (Loading)
2. Extraction (Matching)
3. Expression Building
4. Aggregation

## Input (Batching/Loading)

Input (or batching) is the act of feeding contents read from a file (or stdin/pipe) into
the next stages.  Many times this is invisible, and is simply the pipe or specified filename.

It is possible to tune the batcher to follow the file or batch in different ways.

Read more at:

* [input](input.md)

## Extraction (Matching)

Extraction is denoted with `-m` (match) and is the process of reading a line in
a file or set of files and parsing it with a regular expression into the
match-groups denoted by the regex.

If the regex doesn't match, the line is discarded (a non-match)

These match groups are then fed into the next stage, the expression.

Read more at:

* [extractor](extractor.md)

## Expressions

Expressions `-e` take the match groups, and other pieces of information, and build
a string-based key.  The match groups can be operated on by helpers to build
the string-key (eg arithmatic, json parsing, simple logic).

The result of this key will act as the key for the aggregation stage.

Optional ignore expression(s) can be supplied with `-i` which will
cause the match to be ignored if the expression evaluates to be truthy.

Read more at:

* [expressions](expressions.md)

## Aggregation

The last stage is the aggregation and realtime presentation.  This takes the
key built in the expression stage and uses it to aggregate on. Different aggregators
display the data in different ways.

Aggregator types:

* `filter` is grep-like, in that each line will be processed and the extracted key will be output directly to stdout
* `histogram` will count instances of the extracted key
* `table` will count the key in 2 dimensions
* `bargraph` will create either a stacked or non-stacked bargraph based on 2 dimensions
* `analyze` will use the key as a numeric value and compute mean/median/mode/stddev/percentiles

For more details, see [aggregators](aggregators.md), run `rare --help` for usage,
or look at some [examples](examples.md)
