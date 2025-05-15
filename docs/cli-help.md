## NAME

rare - A fast regex parser, extractor and realtime aggregator

## SYNOPSIS

rare

```
[--color]
[--funcs]=[value]
[--help|-h]
[--metrics]
[--nocolor|--nc]
[--noformat|--nf]
[--noload|--nl]
[--notrim]
[--nounicode|--nu]
[--profile]=[value]
[--version|-v]
```

## DESCRIPTION

Aggregate and display information parsed from text files using
	regex and a simple handlebars-like expressions.

	Run "rare docs overview" or go to https://rare.zdyn.net for more information
	
	https://github.com/zix99/rare

**Usage**:

```
rare [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

## GLOBAL OPTIONS

**--color**: Force-enable color output

**--funcs**="": Specify filenames to load expressions from

**--help, -h**: show help

**--metrics**: Outputs runtime memory metrics after a program runs

**--nocolor, --nc**: Disables color output

**--noformat, --nf**: Disable number formatting

**--noload, --nl**: Disable external file loading in expressions

**--notrim**: By default, rare will trim output text for in-place updates. Setting this flag will disable that

**--nounicode, --nu**: Disable usage of unicode characters

**--profile**="": Write application profiling information as part of execution. Specify base-name

**--version, -v**: print the version


## COMMANDS

### filter, f

Filter incoming results with search criteria, and output raw matches

**--batch**="": Specifies io batching size. Set to 1 for immediate input (default: 1000)

**--batch-buffer**="": Specifies how many batches to read-ahead. Impacts memory usage, can improve performance (default: 16)

**--dissect, -d**="": Dissect expression create match groups to summarize on

**--exclude**="": Glob file patterns to exclude (eg. *.txt)

**--exclude-dir**="": Glob file patterns to exclude directories

**--extract, -e**="": Expression that will generate the key to group by. Specify multiple times for multi-dimensions or use {$} helper (default: [{0}])

**--follow, -f**: Read appended data as file grows

**--follow-symlinks, -L**: Follow symbolic directory links

**--gunzip, -z**: Attempt to decompress file when reading

**--ignore, -i**="": Ignore a match given a truthy expression (Can have multiple)

**--ignore-case, -I**: Augment matcher to be case insensitive

**--include**="": Glob file patterns to include (eg. *.txt)

**--line, -l**: Output source file and line number

**--match, -m**="": Regex to create match groups to summarize on (default: .*)

**--mount**: Don't descend directories on other filesystems

**--num, -n**="": Print the first NUM of lines seen (Not necessarily in-order) (default: 0)

**--poll**: When following a file, poll for changes rather than using inotify

**--posix, -p**: Compile regex as against posix standard

**--read-symlinks**: Read files that are symbolic links

**--readers, --wr**="": Sets the number of concurrent readers (Infinite when -f) (default: 3)

**--recursive, -R**: Recursively walk a non-globbing path and search for plain-files

**--reopen, -F**: Same as -f, but will reopen recreated files

**--tail, -t**: When following a file, navigate to the end of the file to skip existing content

**--workers, -w**="": Set number of data processors (default: 8)

### histogram, histo, h

Summarize results by extracting them to a histogram

**--all, -a**: After summarization is complete, print all histogram buckets

**--atleast**="": Only show results if there are at least this many samples (default: 0)

**--bars, -b**: Display bars as part of histogram

**--batch**="": Specifies io batching size. Set to 1 for immediate input (default: 1000)

**--batch-buffer**="": Specifies how many batches to read-ahead. Impacts memory usage, can improve performance (default: 16)

**--csv, -o**="": Write final results to csv. Use - to output to stdout

**--dissect, -d**="": Dissect expression create match groups to summarize on

**--exclude**="": Glob file patterns to exclude (eg. *.txt)

**--exclude-dir**="": Glob file patterns to exclude directories

**--extra, -x**: Alias for -b --percentage

**--extract, -e**="": Expression that will generate the key to group by. Specify multiple times for multi-dimensions or use {$} helper (default: [{0}])

**--follow, -f**: Read appended data as file grows

**--follow-symlinks, -L**: Follow symbolic directory links

**--format, --fmt**="": Defines a format expression for displayed values

**--gunzip, -z**: Attempt to decompress file when reading

**--ignore, -i**="": Ignore a match given a truthy expression (Can have multiple)

**--ignore-case, -I**: Augment matcher to be case insensitive

**--include**="": Glob file patterns to include (eg. *.txt)

**--match, -m**="": Regex to create match groups to summarize on (default: .*)

**--mount**: Don't descend directories on other filesystems

**--noout**: Don't output any aggregation to stdout

**--num, -n**="": Number of elements to display (default: 5)

**--percentage**: Display percentage of total next to the value

**--poll**: When following a file, poll for changes rather than using inotify

**--posix, -p**: Compile regex as against posix standard

**--read-symlinks**: Read files that are symbolic links

**--readers, --wr**="": Sets the number of concurrent readers (Infinite when -f) (default: 3)

**--recursive, -R**: Recursively walk a non-globbing path and search for plain-files

**--reopen, -F**: Same as -f, but will reopen recreated files

**--scale**="": Defines data-scaling (linear, log10, log2) (default: linear)

**--snapshot**: In aggregators that support it, only output final results, and not progressive updates. Will enable automatically when piping output

**--sort**="": Sorting method for display (value, text, numeric, contextual, date) (default: value)

**--tail, -t**: When following a file, navigate to the end of the file to skip existing content

**--workers, -w**="": Set number of data processors (default: 8)

### heatmap, heat, hm

Create a 2D heatmap of extracted data

**--batch**="": Specifies io batching size. Set to 1 for immediate input (default: 1000)

**--batch-buffer**="": Specifies how many batches to read-ahead. Impacts memory usage, can improve performance (default: 16)

**--cols**="": Number of columns to display (default: 65)

**--csv, -o**="": Write final results to csv. Use - to output to stdout

**--delim**="": Character to tabulate on. Use {$} helper by default (default: \x00)

**--dissect, -d**="": Dissect expression create match groups to summarize on

**--exclude**="": Glob file patterns to exclude (eg. *.txt)

**--exclude-dir**="": Glob file patterns to exclude directories

**--extract, -e**="": Expression that will generate the key to group by. Specify multiple times for multi-dimensions or use {$} helper (default: [{0}])

**--follow, -f**: Read appended data as file grows

**--follow-symlinks, -L**: Follow symbolic directory links

**--format, --fmt**="": Defines a format expression for displayed values

**--gunzip, -z**: Attempt to decompress file when reading

**--ignore, -i**="": Ignore a match given a truthy expression (Can have multiple)

**--ignore-case, -I**: Augment matcher to be case insensitive

**--include**="": Glob file patterns to include (eg. *.txt)

**--match, -m**="": Regex to create match groups to summarize on (default: .*)

**--max**="": Sets the upper bounds of the heatmap (default: auto) (default: 0)

**--min**="": Sets the lower bounds of the heatmap (default: auto) (default: 0)

**--mount**: Don't descend directories on other filesystems

**--noout**: Don't output any aggregation to stdout

**--num, --rows, -n**="": Number of elements (rows) to display (default: 20)

**--poll**: When following a file, poll for changes rather than using inotify

**--posix, -p**: Compile regex as against posix standard

**--read-symlinks**: Read files that are symbolic links

**--readers, --wr**="": Sets the number of concurrent readers (Infinite when -f) (default: 3)

**--recursive, -R**: Recursively walk a non-globbing path and search for plain-files

**--reopen, -F**: Same as -f, but will reopen recreated files

**--scale**="": Defines data-scaling (linear, log10, log2) (default: linear)

**--snapshot**: In aggregators that support it, only output final results, and not progressive updates. Will enable automatically when piping output

**--sort-cols**="": Sorting method for display (value, text, numeric, contextual, date) (default: numeric)

**--sort-rows**="": Sorting method for display (value, text, numeric, contextual, date) (default: numeric)

**--tail, -t**: When following a file, navigate to the end of the file to skip existing content

**--workers, -w**="": Set number of data processors (default: 8)

### spark, sparkline, s

Create rows of sparkline graphs

**--batch**="": Specifies io batching size. Set to 1 for immediate input (default: 1000)

**--batch-buffer**="": Specifies how many batches to read-ahead. Impacts memory usage, can improve performance (default: 16)

**--cols**="": Number of columns to display (default: 65)

**--csv, -o**="": Write final results to csv. Use - to output to stdout

**--delim**="": Character to tabulate on. Use {$} helper by default (default: \x00)

**--dissect, -d**="": Dissect expression create match groups to summarize on

**--exclude**="": Glob file patterns to exclude (eg. *.txt)

**--exclude-dir**="": Glob file patterns to exclude directories

**--extract, -e**="": Expression that will generate the key to group by. Specify multiple times for multi-dimensions or use {$} helper (default: [{0}])

**--follow, -f**: Read appended data as file grows

**--follow-symlinks, -L**: Follow symbolic directory links

**--format, --fmt**="": Defines a format expression for displayed values

**--gunzip, -z**: Attempt to decompress file when reading

**--ignore, -i**="": Ignore a match given a truthy expression (Can have multiple)

**--ignore-case, -I**: Augment matcher to be case insensitive

**--include**="": Glob file patterns to include (eg. *.txt)

**--match, -m**="": Regex to create match groups to summarize on (default: .*)

**--mount**: Don't descend directories on other filesystems

**--noout**: Don't output any aggregation to stdout

**--notruncate**: Disable truncating data that doesn't fit in the sparkline

**--num, --rows, -n**="": Number of elements (rows) to display (default: 20)

**--poll**: When following a file, poll for changes rather than using inotify

**--posix, -p**: Compile regex as against posix standard

**--read-symlinks**: Read files that are symbolic links

**--readers, --wr**="": Sets the number of concurrent readers (Infinite when -f) (default: 3)

**--recursive, -R**: Recursively walk a non-globbing path and search for plain-files

**--reopen, -F**: Same as -f, but will reopen recreated files

**--scale**="": Defines data-scaling (linear, log10, log2) (default: linear)

**--snapshot**: In aggregators that support it, only output final results, and not progressive updates. Will enable automatically when piping output

**--sort-cols**="": Sorting method for display (value, text, numeric, contextual, date) (default: numeric)

**--sort-rows**="": Sorting method for display (value, text, numeric, contextual, date) (default: value)

**--tail, -t**: When following a file, navigate to the end of the file to skip existing content

**--workers, -w**="": Set number of data processors (default: 8)

### bargraph, bars, bar, b

Create a bargraph of the given 1 or 2 dimension data

**--batch**="": Specifies io batching size. Set to 1 for immediate input (default: 1000)

**--batch-buffer**="": Specifies how many batches to read-ahead. Impacts memory usage, can improve performance (default: 16)

**--csv, -o**="": Write final results to csv. Use - to output to stdout

**--dissect, -d**="": Dissect expression create match groups to summarize on

**--exclude**="": Glob file patterns to exclude (eg. *.txt)

**--exclude-dir**="": Glob file patterns to exclude directories

**--extract, -e**="": Expression that will generate the key to group by. Specify multiple times for multi-dimensions or use {$} helper (default: [{0}])

**--follow, -f**: Read appended data as file grows

**--follow-symlinks, -L**: Follow symbolic directory links

**--format, --fmt**="": Defines a format expression for displayed values

**--gunzip, -z**: Attempt to decompress file when reading

**--ignore, -i**="": Ignore a match given a truthy expression (Can have multiple)

**--ignore-case, -I**: Augment matcher to be case insensitive

**--include**="": Glob file patterns to include (eg. *.txt)

**--match, -m**="": Regex to create match groups to summarize on (default: .*)

**--mount**: Don't descend directories on other filesystems

**--noout**: Don't output any aggregation to stdout

**--poll**: When following a file, poll for changes rather than using inotify

**--posix, -p**: Compile regex as against posix standard

**--read-symlinks**: Read files that are symbolic links

**--readers, --wr**="": Sets the number of concurrent readers (Infinite when -f) (default: 3)

**--recursive, -R**: Recursively walk a non-globbing path and search for plain-files

**--reopen, -F**: Same as -f, but will reopen recreated files

**--scale**="": Defines data-scaling (linear, log10, log2) (default: linear)

**--snapshot**: In aggregators that support it, only output final results, and not progressive updates. Will enable automatically when piping output

**--sort**="": Sorting method for display (value, text, numeric, contextual, date) (default: numeric)

**--stacked, -s**: Display bargraph as stacked

**--tail, -t**: When following a file, navigate to the end of the file to skip existing content

**--workers, -w**="": Set number of data processors (default: 8)

### analyze, a

Numerical analysis on a set of filtered data

**--batch**="": Specifies io batching size. Set to 1 for immediate input (default: 1000)

**--batch-buffer**="": Specifies how many batches to read-ahead. Impacts memory usage, can improve performance (default: 16)

**--dissect, -d**="": Dissect expression create match groups to summarize on

**--exclude**="": Glob file patterns to exclude (eg. *.txt)

**--exclude-dir**="": Glob file patterns to exclude directories

**--extra, -x**: Displays extra analysis on the data (Requires more memory and cpu)

**--extract, -e**="": Expression that will generate the key to group by. Specify multiple times for multi-dimensions or use {$} helper (default: [{0}])

**--follow, -f**: Read appended data as file grows

**--follow-symlinks, -L**: Follow symbolic directory links

**--gunzip, -z**: Attempt to decompress file when reading

**--ignore, -i**="": Ignore a match given a truthy expression (Can have multiple)

**--ignore-case, -I**: Augment matcher to be case insensitive

**--include**="": Glob file patterns to include (eg. *.txt)

**--match, -m**="": Regex to create match groups to summarize on (default: .*)

**--mount**: Don't descend directories on other filesystems

**--poll**: When following a file, poll for changes rather than using inotify

**--posix, -p**: Compile regex as against posix standard

**--quantile, -q**="": Adds a quantile to the output set. Requires --extra (default: [90 99 99.9])

**--read-symlinks**: Read files that are symbolic links

**--readers, --wr**="": Sets the number of concurrent readers (Infinite when -f) (default: 3)

**--recursive, -R**: Recursively walk a non-globbing path and search for plain-files

**--reopen, -F**: Same as -f, but will reopen recreated files

**--reverse, -r**: Reverses the numerical series when ordered-analysis takes place (eg Quantile)

**--snapshot**: In aggregators that support it, only output final results, and not progressive updates. Will enable automatically when piping output

**--tail, -t**: When following a file, navigate to the end of the file to skip existing content

**--workers, -w**="": Set number of data processors (default: 8)

### tabulate, table, t

Create a 2D summarizing table of extracted data

**--batch**="": Specifies io batching size. Set to 1 for immediate input (default: 1000)

**--batch-buffer**="": Specifies how many batches to read-ahead. Impacts memory usage, can improve performance (default: 16)

**--cols**="": Number of columns to display (default: 10)

**--coltotal**: Show column totals

**--csv, -o**="": Write final results to csv. Use - to output to stdout

**--delim**="": Character to tabulate on. Use {$} helper by default (default: \x00)

**--dissect, -d**="": Dissect expression create match groups to summarize on

**--exclude**="": Glob file patterns to exclude (eg. *.txt)

**--exclude-dir**="": Glob file patterns to exclude directories

**--extra, -x**: Display row and column totals

**--extract, -e**="": Expression that will generate the key to group by. Specify multiple times for multi-dimensions or use {$} helper (default: [{0}])

**--follow, -f**: Read appended data as file grows

**--follow-symlinks, -L**: Follow symbolic directory links

**--format, --fmt**="": Defines a format expression for displayed values

**--gunzip, -z**: Attempt to decompress file when reading

**--ignore, -i**="": Ignore a match given a truthy expression (Can have multiple)

**--ignore-case, -I**: Augment matcher to be case insensitive

**--include**="": Glob file patterns to include (eg. *.txt)

**--match, -m**="": Regex to create match groups to summarize on (default: .*)

**--mount**: Don't descend directories on other filesystems

**--noout**: Don't output any aggregation to stdout

**--num, --rows, -n**="": Number of elements to display (default: 20)

**--poll**: When following a file, poll for changes rather than using inotify

**--posix, -p**: Compile regex as against posix standard

**--read-symlinks**: Read files that are symbolic links

**--readers, --wr**="": Sets the number of concurrent readers (Infinite when -f) (default: 3)

**--recursive, -R**: Recursively walk a non-globbing path and search for plain-files

**--reopen, -F**: Same as -f, but will reopen recreated files

**--rowtotal**: Show row totals

**--snapshot**: In aggregators that support it, only output final results, and not progressive updates. Will enable automatically when piping output

**--sort-cols**="": Sorting method for display (value, text, numeric, contextual, date) (default: value)

**--sort-rows**="": Sorting method for display (value, text, numeric, contextual, date) (default: value)

**--tail, -t**: When following a file, navigate to the end of the file to skip existing content

**--workers, -w**="": Set number of data processors (default: 8)

### reduce, r

Aggregate the results of a query based on an expression, pulling customized summary from the extracted data

**--accumulator, -a**="": Specify one or more expressions to execute for each match. `{.}` is the accumulator. Syntax: `[name[:initial]=]expr`

**--batch**="": Specifies io batching size. Set to 1 for immediate input (default: 1000)

**--batch-buffer**="": Specifies how many batches to read-ahead. Impacts memory usage, can improve performance (default: 16)

**--cols**="": Number of columns to display (default: 10)

**--csv, -o**="": Write final results to csv. Use - to output to stdout

**--dissect, -d**="": Dissect expression create match groups to summarize on

**--exclude**="": Glob file patterns to exclude (eg. *.txt)

**--exclude-dir**="": Glob file patterns to exclude directories

**--extract, -e**="": Expression that will generate the key to group by. Specify multiple times for multi-dimensions or use {$} helper (default: [{@}])

**--follow, -f**: Read appended data as file grows

**--follow-symlinks, -L**: Follow symbolic directory links

**--format, --fmt**="": Defines a format expression for displayed values. Syntax: `[name=]expr`

**--group, -g**="": Specifies one or more expressions to group on. Syntax: `[name=]expr`

**--gunzip, -z**: Attempt to decompress file when reading

**--ignore, -i**="": Ignore a match given a truthy expression (Can have multiple)

**--ignore-case, -I**: Augment matcher to be case insensitive

**--include**="": Glob file patterns to include (eg. *.txt)

**--initial**="": Specify the default initial value for any accumulators that don't specify (default: 0)

**--match, -m**="": Regex to create match groups to summarize on (default: .*)

**--mount**: Don't descend directories on other filesystems

**--noout**: Don't output any aggregation to stdout

**--num, --rows, -n**="": Number of elements to display (default: 20)

**--poll**: When following a file, poll for changes rather than using inotify

**--posix, -p**: Compile regex as against posix standard

**--read-symlinks**: Read files that are symbolic links

**--readers, --wr**="": Sets the number of concurrent readers (Infinite when -f) (default: 3)

**--recursive, -R**: Recursively walk a non-globbing path and search for plain-files

**--reopen, -F**: Same as -f, but will reopen recreated files

**--snapshot**: In aggregators that support it, only output final results, and not progressive updates. Will enable automatically when piping output

**--sort**="": Specify an expression to sort groups by. Will sort result in alphanumeric order

**--sort-reverse**: Reverses sort order

**--table**: Force output to be a table, even when there are no groups

**--tail, -t**: When following a file, navigate to the end of the file to skip existing content

**--workers, -w**="": Set number of data processors (default: 8)

### docs

Access detailed documentation

**--no-pager, -n**: Don't use pager to view documentation

### expression, exp, expr

Evaluate and benchmark expressions

**--benchmark, -b**: Benchmark the expression (slow)

**--data, -d**="": Specify positional data in the expression

**--key, -k**="": Specify a named argument, a=b

**--listfuncs**: Lists all available expression functions

**--no-optimize**: Disable expression static analysis optimization

**--raw, -r**: Don't format arrays, output raw with null-separators

**--skip-newline, -n**: Don't add a newline character when printing plain result

**--stats, -s**: Display stats about the expression

### help, h

Shows a list of commands or help for one command
