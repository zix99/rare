# rare

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/zix99/rare/rare)](https://github.com/zix99/rare/actions)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/zix99/rare)](https://github.com/zix99/rare/releases)
[![codecov](https://codecov.io/gh/zix99/rare/branch/master/graph/badge.svg)](https://codecov.io/gh/zix99/rare)
![GitHub all releases](https://img.shields.io/github/downloads/zix99/rare/total)
![GitHub](https://img.shields.io/github/license/zix99/rare)

A fast text scanner/regex extractor and realtime summarizer.

Supports various CLI-based graphing and metric formats (filter (grep-like), histogram, table, bargraph, etc).

`rare` is a play on "more" and "less", but can also stand for "realtime aggregated regular expressions".

See [rare.zdyn.net](https://rare.zdyn.net) or the [docs/ folder](docs/) for the full documentation.

![rare gif](docs/images/rare.gif)

## Features

 * Multiple summary formats including: filter (like grep), histogram, bar graphs, tables, heatmaps, and numerical analysis
 * File glob expansions (eg `/var/log/*` or `/var/log/*/*.log`) and `-R`
 * Optional gzip decompression (with `-z`)
 * Following `-f` or re-open following `-F` (use `--poll` to poll, and `--tail` to tail)
 * Ignoring lines that match an expression (with `-i`)
 * Aggregating and realtime summary (Don't have to wait for all data to be scanned)
 * Multi-threaded reading, parsing, and aggregation (It's fast)
 * Color-coded outputs (optionally)
 * Pipe support (stdin for reading, stdout will disable color) eg. `tail -f | rare ...`

Take a look at [examples](docs/usage/examples.md) to see more of what *rare* does.

### Output Formats

Output formats include:

* `filter` is grep-like, in that each line will be processed and the extracted key will be output directly to stdout
* `histogram` will count instances of the extracted key
* `table` will count the key in 2 dimensions
* `heatmap` will display a color-coded version of the strength of a cell in a dense format
* `bargraph` will create either a stacked or non-stacked bargraph based on 2 dimensions
* `analyze` will use the key as a numeric value and compute mean/median/mode/stddev/percentiles

More details on various output formats and aggregators (including examples) can be found in [aggregators](docs/usage/aggregators.md)

## Installation

### Manual

Download appropriate binary or package from [Releases](https://github.com/zix99/rare/releases)

### Homebrew

```sh
brew tap zix99/rare
brew install rare
```

### Community Contributed

The below install methods have been contributed by the community, and aren't maintained directly.

#### MacPorts

```sh
sudo port selfupdate
sudo port install rare
```

### From code

Clone the repo, and:

Requires GO 1.17 or higher

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

**A Note on PCRE (Perl Compatible Regex Library)**

Besides your standard OS versions, there is an additional `pcre` build which is ~4x faster than go's `re2` implementation in moderately complex cases.  In order to use this, you must make sure that **libpcre2** is installed (eg `apt install libpcre2-8-0`).  Right now, it is only bundled with the linux distribution.

PCRE2 also comes with pitfalls, two of the most important are:
1. That *rare* is now dynamically linked, meaning that you need to have libc and libpcre installed
2. That pcre is an exponential-time algorithm (re2 is linear).  While it can be significantly faster than go's `re2`, it can also be catastropically slower in some situations. There is a good post [here](https://swtch.com/~rsc/regexp/regexp1.html) that talks about regexp timings.

I will leave it up to the user as to which they find suitable to use for their situation.  Generally, if you know what *rare* is getting as an input, the pcre version is perfectly safe and can be much faster.

## Documentation

All documentation may be found here, in the [docs/](docs/) folder, by running `rare docs` (embedded docs/ folder), or on the website [rare.zdyn.net](https://rare.zdyn.net)

You can also see a dump of the CLI options at [cli-help.md](docs/cli-help.md)

## Example

### Create histogram from sample data

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

### Extact status and size from nginx logs
```sh
$ rare filter -n 4 -m "(\d{3}) (\d+)" -e "{1} {2}" access.log
404 169
404 169
404 571
404 571
Matched: 4 / 4
```

### Extract status codes from nginx logs

```sh
$ rare histo \
    -m '"(\w{3,4}) ([A-Za-z0-9/.]+).*" (\d{3})' \ # The regex that extracts match-groups
    -e '{3} {1}' \ # The expression will be the key, referencing the match-groups
    access.log     # One or more files (or -R for recursion)

200 GET                          160663
404 GET                          857
304 GET                          53
200 HEAD                         18
403 GET                          14
```

### More Examples

For more examples, check out the [docs](docs/usage/examples.md) or [the website](https://rare.zdyn.net/usage/examples/)


## Performance Benchmarking

I know there are different solutions, and rare accomplishes summarization in a way
that grep, awk, etc can't, however I think it's worth analyzing the performance of this
tool vs standard tools to show that it's at least as good.

See [benchmarks](docs/benchmarks.md) for comparisons between common tools like `grep | wc`,
silversearcher-ag, etc.


## Development

New additions to `rare` should pass the following checks

- Documentation for any new functionality or expression changes
- Before and after CPU and memory benchmarking for core additions (Expressions, aggregation, benchmarking, and rendering)
- Limit memory allocations (preferably 0!) in the high-throughput functions
- Tests, and if it makes sense, benchmarks of a given function

### Running/Testing

```bash
go run .
go test ./...
```

### Profiling

New high-throughput changes should be performance benchmarked.

To Benchmark:

```bash
go run . --profile out <your test code>
go tool pprof -http=:8080 out.cpu.prof # CPU
go tool pprof -http=:8080 out_num.prof # Memory
```

## License

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
