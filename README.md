# rare

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/zix99/rare/rare) ![GitHub release (latest by date)](https://img.shields.io/github/v/release/zix99/rare) [![codecov](https://codecov.io/gh/zix99/rare/branch/master/graph/badge.svg)](https://codecov.io/gh/zix99/rare) ![GitHub all releases](https://img.shields.io/github/downloads/zix99/rare/total)

A fast text scanner/regex extractor and realtime summarizer.

Supports various CLI-based graphing and metric formats (filter (grep-like), histogram, table, bargraph, etc).

`rare` is a play on "more" and "less", but can also stand for "realtime aggregated regular expressions".

![rare gif](images/rare.gif)

# Features

 * Multiple summary formats including: filter (like grep), histogram, bar graphs, and numerical analysis
 * File glob expansions (eg `/var/log/*` or `/var/log/*/*.log`) and `-R`
 * Optional gzip decompression (with `-z`)
 * Following `-f` or re-open following `-F` (use `--poll` to poll)
 * Ignoring lines that match an expression (with `-i`)
 * Aggregating and realtime summary (Don't have to wait for all data to be scanned)
 * Multi-threaded reading, parsing, and aggregation
 * Color-coded outputs (optionally)
 * Pipe support (stdin for reading, stdout will disable color) eg. `tail -f | rare ...`

# Installation

**A Note on PCRE**

Besides your standard OS versions, there is an additional `pcre` build which is 4x faster than go's `re2` implementation in moderately complex cases.  In order to use this, you must make sure that libpcre2 is installed (eg `apt install libpcre2-8-0`).  Right now, it is only bundled with the linux distribution.

PCRE2 also comes with pitfalls, two of the most important are:
1. That *rare* is now dynamically linked, meaning that you need to have libc and libpcre installed
2. That pcre is an exponential-time algorithm (re2 is linear).  While it can be significantly faster than go's `re2`, it can also be catastropically slower in some situations. There is a good post [here](https://swtch.com/~rsc/regexp/regexp1.html) that talks about regexp timings.

I will leave it up to the user as to which they find suitable to use for their situation.  Generally, if you know what *rare* is getting as an input, the pcre version is perfectly safe and can be much faster.

## Manual

Download appropriate binary or package from [Releases](https://github.com/zix99/rare/releases)

## Homebrew

```sh
brew tap zix99/rare
brew install rare
```

## From code

Clone the repo, and:

Requires GO 1.16 or higher

```sh
go mod download

# Build binary
go build .

# OR, with experimental features
go build -tags experimental .
```

Available tags:

* `experimental` Enable experimental features (eg. fuzzy search)
* `pcre2` Enables PCRE 2 (v10) where able. Currently linux only

# Documentation

All documentation may be found here, in the [docs/](docs/) folder, and by running `rare docs` (embedded docs/ folder)

You can also see a dump of the CLI options at [cli-help.md](cli-help.md)

# Example

## Create histogram from sample data

```sh
$ cat input.txt
1
2
1
3
1
0

$ rare histo input.txt
1                   3         
0                   1         
2                   1         
3                   1         

Matched: 6 / 6 (Groups: 4)
```

## Extact status and size from nginx logs
```sh
$ rare filter -n 4 -m "(\d{3}) (\d+)" -e "{1} {2}" access.log
404 169
404 169
404 571
404 571
Matched: 4 / 4
```

## Extract status codes from nginx logs

```sh
$ rare histo -m '"(\w{3,4}) ([A-Za-z0-9/.]+).*" (\d{3})' -e '{3} {1}' access.log
200 GET                          160663
404 GET                          857
304 GET                          53
200 HEAD                         18
403 GET                          14
```

## Extract number of bytes sent by bucket, and format

This shows an example of how to bucket the values into size of `1000`. In this case, it doesn't make
sense to see the histogram by number of bytes, but we might want to know the ratio of various orders-of-magnitudes.

```sh
$ rare histo -m '"(\w{3,4}) ([A-Za-z0-9/.]+).*" (\d{3}) (\d+)' -e "{bucket {4} 10000}" -n 10 access.log -b
0                   144239     ||||||||||||||||||||||||||||||||||||||||||||||||||
190000              2599       
10000               1290       
180000              821        
20000               496        
30000               445        
40000               440        
200000              427        
140000              323        
70000               222        
Matched: 161622 / 161622
Groups:  1203
```

# Output Formats

## Histogram (histo)

The histogram format outputs an aggregation by counting the occurences of an extracted match.  That is to say, on every line a regex will be matched (or not), and the matched groups can be used to extract and build a key, that will act as the bucketing name.

```
NAME:
   rare histogram - Summarize results by extracting them to a histogram

USAGE:
   rare histogram [command options] <-|filename|glob...>

DESCRIPTION:
   Generates a live-updating histogram of the extracted information from a file
    Each line in the file will be matched, any the matching part extracted
    as a key and counted.
    If an extraction expression is provided with -e, that will be used
    as the key instead

OPTIONS:
   --follow, -f                 Read appended data as file grows
   --reopen, -F                 Same as -f, but will reopen recreated files
   --poll                       When following a file, poll for changes rather than using inotify
   --posix, -p                  Compile regex as against posix standard
   --match value, -m value      Regex to create match groups to summarize on (default: ".*")
   --extract value, -e value    Expression that will generate the key to group by (default: "{0}")
   --gunzip, -z                 Attempt to decompress file when reading
   --batch value                Specifies io batching size. Set to 1 for immediate input (default: 1000)
   --workers value, -w value    Set number of data processors (default: 5)
   --readers value, --wr value  Sets the number of concurrent readers (Infinite when -f) (default: 3)
   --ignore value, -i value     Ignore a match given a truthy expression (Can have multiple)
   --recursive, -R              Recursively walk a non-globbing path and search for plain-files
   --bars, -b                   Display bars as part of histogram
   --num value, -n value        Number of elements to display (default: 5)
   --reverse                    Reverses the display sort-order
   --sortkey, --sk              Sort by key, rather than value

```

## Filter (filter)

Filter is a command used to match and (optionally) extract that match without any aggregation. It's effectively a `grep` or a combination of `grep`, `awk`, and/or `sed`.

```
NAME:
   rare filter - Filter incoming results with search criteria, and output raw matches

USAGE:
   rare filter [command options] <-|filename|glob...>

DESCRIPTION:
   Filters incoming results by a regex, and output the match or an extracted expression.
    Unable to output contextual information due to the application's parallelism.  Use grep if you
    need that

OPTIONS:
   --follow, -f                 Read appended data as file grows
   --reopen, -F                 Same as -f, but will reopen recreated files
   --poll                       When following a file, poll for changes rather than using inotify
   --posix, -p                  Compile regex as against posix standard
   --match value, -m value      Regex to create match groups to summarize on (default: ".*")
   --extract value, -e value    Expression that will generate the key to group by (default: "{0}")
   --gunzip, -z                 Attempt to decompress file when reading
   --batch value                Specifies io batching size. Set to 1 for immediate input (default: 1000)
   --workers value, -w value    Set number of data processors (default: 5)
   --readers value, --wr value  Sets the number of concurrent readers (Infinite when -f) (default: 3)
   --ignore value, -i value     Ignore a match given a truthy expression (Can have multiple)
   --recursive, -R              Recursively walk a non-globbing path and search for plain-files
   --line, -l                   Output line numbers
```

## Bar Graph

Similar to histogram or table, bargraph can generate a stacked or grouped bargraph by one or two keys.

```sh
$ rare bars -sz -m "\[(.+?)\].*\" (\d+)" -e "{$ {buckettime {1} year nginx} {bucket {2} {multi 10 10}}}" testdata/*

        █ 200  █ 400  █ 300
2019  ████████████████████████████████████████  3,741,444
2020  █████████████████████████████████████████████████  4,631,884
Matched: 8,373,328 / 8,383,717
```

```
NAME:
   rare bargraph - Create a bargraph of the given 1 or 2 dimension data

USAGE:
   rare bargraph [command options] <-|filename|glob...>

DESCRIPTION:
   Creates a bargraph of one or two dimensional data.  Unlike histogram
    the bargraph can collapse and stack data in different formats.  The key data format
    is {$ a b [c]}, where a is the base-key, b is the optional sub-key, and c is the increment
    (defeaults to 1)

OPTIONS:
   --follow, -f                 Read appended data as file grows
   --reopen, -F                 Same as -f, but will reopen recreated files
   --poll                       When following a file, poll for changes rather than using inotify
   --posix, -p                  Compile regex as against posix standard
   --match value, -m value      Regex to create match groups to summarize on (default: ".*")
   --extract value, -e value    Expression that will generate the key to group by (default: "{0}")
   --gunzip, -z                 Attempt to decompress file when reading
   --batch value                Specifies io batching size. Set to 1 for immediate input (default: 1000)
   --workers value, -w value    Set number of data processors (default: 3)
   --readers value, --wr value  Sets the number of concurrent readers (Infinite when -f) (default: 3)
   --ignore value, -i value     Ignore a match given a truthy expression (Can have multiple)
   --recursive, -R              Recursively walk a non-globbing path and search for plain-files
   --stacked, -s                Display bargraph as stacked
   --reverse                    Reverses the display sort-order
```

## Numerical Analysis

This command will extract a number from logs and run basic analysis on that number (Such as mean, median, mode, and quantiles).

```
NAME:
   rare analyze - Numerical analysis on a set of filtered data

USAGE:
   rare analyze [command options] <-|filename|glob...>

DESCRIPTION:
   Treat every extracted expression as a numerical input, and run analysis
    on that input.  Will extract mean, median, mode, min, max.  If specifying --extra
    will also extract std deviation, and quantiles

OPTIONS:
   --follow, -f                 Read appended data as file grows
   --reopen, -F                 Same as -f, but will reopen recreated files
   --poll                       When following a file, poll for changes rather than using inotify
   --posix, -p                  Compile regex as against posix standard
   --match value, -m value      Regex to create match groups to summarize on (default: ".*")
   --extract value, -e value    Expression that will generate the key to group by (default: "{0}")
   --gunzip, -z                 Attempt to decompress file when reading
   --batch value                Specifies io batching size. Set to 1 for immediate input (default: 1000)
   --workers value, -w value    Set number of data processors (default: 5)
   --readers value, --wr value  Sets the number of concurrent readers (Infinite when -f) (default: 3)
   --ignore value, -i value     Ignore a match given a truthy expression (Can have multiple)
   --recursive, -R              Recursively walk a non-globbing path and search for plain-files
   --extra                      Displays extra analysis on the data (Requires more memory and cpu)
   --reverse, -r                Reverses the numerical series when ordered-analysis takes place (eg Quantile)
   --quantile value, -q value   Adds a quantile to the output set. Requires --extra (default: "90", "99", "99.9")
```

**Example:**

```bash
$ go run *.go --color analyze -m '"(\w{3,4}) ([A-Za-z0-9/.@_-]+).*" (\d{3}) (\d+)' -e "{4}" testdata/access.log 
Samples:  161,622
Mean:     2,566,283.9616
Min:      0.0000
Max:      1,198,677,592.0000

Median:   1,021.0000
Mode:     1,021.0000
P90:      19,506.0000
P99:      64,757,808.0000
P99.9:    395,186,166.0000
Matched: 161,622 / 161,622
```

## Tabulate

Create a 2D view (table) of data extracted from a file. Expression needs to yield a two dimensions separated by a tab.  Can either use `\x00` or the `{$ a b}` helper.  First element is the column name, followed by the row name.

```
NAME:
   rare tabulate - Create a 2D summarizing table of extracted data

USAGE:
   rare tabulate [command options] <-|filename|glob...>

DESCRIPTION:
   Summarizes the extracted data as a 2D data table.
    The key is provided in the expression, and should be separated by a tab \x00
    character or via {$ a b} Where a is the column header, and b is the row

OPTIONS:
   --follow, -f                 Read appended data as file grows
   --reopen, -F                 Same as -f, but will reopen recreated files
   --poll                       When following a file, poll for changes rather than using inotify
   --posix, -p                  Compile regex as against posix standard
   --match value, -m value      Regex to create match groups to summarize on (default: ".*")
   --extract value, -e value    Expression that will generate the key to group by (default: "{0}")
   --gunzip, -z                 Attempt to decompress file when reading
   --batch value                Specifies io batching size. Set to 1 for immediate input (default: 1000)
   --workers value, -w value    Set number of data processors (default: 5)
   --readers value, --wr value  Sets the number of concurrent readers (Infinite when -f) (default: 3)
   --ignore value, -i value     Ignore a match given a truthy expression (Can have multiple)
   --recursive, -R              Recursively walk a non-globbing path and search for plain-files
   --delim value                Character to tabulate on. Use {$} helper by default (default: "\x00")
   --num value, -n value        Number of elements to display (default: 20)
   --cols value                 Number of columns to display (default: 10)
   --sortkey, --sk              Sort rows by key name rather than by values
```

**Example:**

```bash
$ rare tabulate -m "(\d{3}) (\d+)" -e "{$ {1} {bucket {2} 100000}}" -sk access.log

         200      404      304      403      301      206      
0        153,271  860      53       14       12       2                 
1000000  796      0        0        0        0        0                 
2000000  513      0        0        0        0        0                 
7000000  262      0        0        0        0        0                 
4000000  257      0        0        0        0        0                 
6000000  221      0        0        0        0        0                 
5000000  218      0        0        0        0        0                 
9000000  206      0        0        0        0        0                 
3000000  202      0        0        0        0        0                 
10000000 201      0        0        0        0        0                 
11000000 190      0        0        0        0        0                 
21000000 142      0        0        0        0        0                 
15000000 138      0        0        0        0        0                 
8000000  137      0        0        0        0        0                 
22000000 123      0        0        0        0        0                 
14000000 121      0        0        0        0        0                 
16000000 110      0        0        0        0        0                 
17000000 99       0        0        0        0        0                 
34000000 91       0        0        0        0        0                 
Matched: 161,622 / 161,622
Rows: 223; Cols: 6
```

# Performance Benchmarking

I know there are different solutions, and rare accomplishes summarization in a way
that grep, awk, etc can't, however I think it's worth analyzing the performance of this
tool vs standard tools to show that it's at least as good.

It's worth noting that in many of these results rare is just as fast, but part
of that reason is that it consumes CPU in a more efficient way (go is great at parallelization).
So take that into account, for better or worse.

All tests were done on ~83MB of gzip'd (1.5GB gunzip'd) nginx logs spread across 10 files.

Each program was run 3 times and the last time was taken (to make sure things were cached equally).

## zcat & grep

```
$ time zcat testdata/* | grep -Poa '" (\d{3})' | wc -l
8373328

real    0m11.272s
user    0m16.239s
sys     0m1.989s

$ time zcat testdata/* | grep -Poa '" 200' > /dev/null

real    0m5.416s
user    0m4.810s
sys     0m1.185s
```

I believe the largest holdup here is the fact that zcat will pass all the data to grep via a synchronous pipe, whereas
rare can process everything in async batches.  Using `pigz` instead didn't yield different results, but on single-file
results they did perform comparibly.

## Silver Searcher (ag)

ag version 2.2.0 has a bug where it won't scan all my testdata.  I'll hold on benchmarking until there's a fix.

### Old Benchmark (Less data by factor of ~8x)
```
$ ag --version
ag version 2.2.0

Features:
  +jit +lzma +zlib

$ time ag -z '" (\d{3})' testdata/* | wc -l
1131354

real	0m3.944s
user	0m3.904s
sys	0m0.152s
```

## rare

At no point scanning the data does `rare` exceed ~76MB of resident memory.

```
$ rare -v
rare version 0.1.16, 11ca2bfc4ad35683c59929a74ad023cc762a29ae

$ time rare filter -m '" (\d{3})' -e "{1}" -z testdata/* | wc -l
Matched: 8,373,328 / 8,373,328
8373328

real    0m16.192s
user    0m20.298s
sys     0m20.697s

$ time rare histo -m '" (\d{3})' -e "{1}" -z testdata/*
404                 5,557,374 
200                 2,564,984 
400                 243,282   
405                 5,708     
408                 1,397     
Matched: 8,373,328 / 8,373,328 (Groups: 8)


real    0m3.869s
user    0m13.423s
sys     0m0.191s
```

### pcre2

The PCRE2 version is approximately the same on a simple regular expression, but begins to shine
on more complex regex's.

```
$ time rare table -z -m "\[(.+?)\].*\" (\d+)" -e "{buckettime {1} year nginx}" -e "{bucket {2} 100}" testdata/*
          2020      2019      
400       2,915,487 2,892,274           
200       1,716,107 848,925             
300       290       245                 
Matched: 8,373,328 / 8,373,328 (R: 3; C: 2)


real    0m31.419s
user    1m40.060s
sys     0m0.657s

$ time rare-pcre table -z -m "\[(.+?)\].*\" (\d+)" -e "{buckettime {1} year nginx}" -e "{bucket {2} 100}" testdata/*
          2020      2019      
400       2,915,487 2,892,274           
200       1,716,107 848,925             
300       290       245                 
Matched: 8,373,328 / 8,373,328 (R: 3; C: 2)


real    0m7.936s
user    0m27.600s
sys     0m0.301s
```

# Development

New additions to `rare` should pass the following checks

- Documentation for any new functionality or expression changes
- Before and after CPU and memory benchmarking for core additions (Expressions, aggregation, benchmarking, and rendering)
- Limit memory allocations (preferably 0!) in the high-throughput functions
- Tests, and if it makes sense, benchmarks of a given function

## Running/Testing

```bash
go run .
go test ./...
```

## Profiling

New high-throughput changes should be performance benchmarked.

To Benchmark:

```bash
go run . --profile out <your test code>
go tool pprof -http=:8080 out.cpu.prof # CPU
go tool pprof -http=:8080 out_num.prof # Memory
```

# License

    Copyright (C) 2019  Christopher LaPointe

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
