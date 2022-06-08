# Input

*rare* reads the supplied inputs in massive parallelization, rather
than in-order reads.  In most cases, you won't need to do anything
other than specifying what to read.  In some cases, you may want to
tweak some parameters.

## Input Methods

### Read File(s)

The simplest version of reading files is by specifying one or more filename:

`rare <aggregator> file1 file2 file3...`

You can also use simple expansions, such as:

`rare <aggregator> path/**/*.log`

In this case, all `*.log` files in any nested directory under `path/` will be read.

or you can use recursion, which will read all plain files in the path

`rare <aggregator> -R path/`

#### gzip

If the files *may* be gzip'd you can specify `-z`, and will be gunzip'd if able.  If a
file can't be opened as a gzip file, a warning will be logged, and it will be interpreted
as a raw file.

`rare <aggregator> -z *.log.gz`

### Following File(s)

Like `tail -f`, following files allows you to watch files actively being written to. This is
useful, for example, to read a log of an actively running application.

**Note:** When following files, all files are open at once, and max readers are ignored.

`rare <aggregator> -f app.log`

If the file may be deleted and recreated, such as in a log-rotation, you can follow with re-open

`rare <aggregator> -F app.log`

#### Polling (Instead of blocking)

By default, following a file uses `fsnotify` which monitors files for changes.  This should
work fine for most major operating systems.  If not, you can enable polling to watch for changes
instead with `--poll`

#### Tailing

If you wish to only start reading at the end of the file (eg. only looking at newer entries),
you can specify `-t` or `--tail` to start following at the end.

### Stdin/Pipe

There are two ways to read from a pipe: implicit and explicit.

Implicitely, if *rare* detects its stdin is a pipe, it will read it simply by not providing any arguments

`cat file.log | rare <aggregator>` or `rare <aggregator> < file.log`

Explicitely, you can pass a single read argument of `-` (dash) to mandate reading from stdin

`cat file.log | rare <aggregator> -`

## Tweaking the Batcher

There are already some heuristics that optimize how files are read which
should work for most cases. If you do find you need to modify how *rare*
is reading, you can tweak two things:

* concurrency -- How many files are read at once
* batch size -- How many lines read from a given file are "batched" to send to the expression stage

### Concurrency

Concurrency specifies how many files are opened at once (in a normal case). It
defaults to `3`, but is ignored if following files.

Specify with:

`rare <aggregator> --readers=1 file1 file2 file3...`

### Batch Sizes

Rare reads (by default) 1000 lines in a file, for a batch, before providing it
to the extractor stage.  This significantly speeds up processing, but comes
at the cost of being less real-time if input generation is slow.

To counteract this, in the *follow* or *stdin* cases, there's also a flush timeout of
250ms. This means if a new line has been received, and the duration has passed,
that the batch will be processed irregardless of its current size.

You can tweak this value with `--batch`

`rare <aggreagator> --batch=10 ...`
