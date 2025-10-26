# Benchmarks

Below, rare is compared to various other common and popular tools.

It's worth noting that in many of these results rare is just as fast, but part
of that reason is that it consumes CPU in a more efficient way (go is great at parallelization).
So take that into account, for better or worse.

All tests were done on ~824MB of gzip'd (13.93 GB gunzip'd) nginx logs spread across 8 files.  They
were run on a NVMe SSD on a recent (2025) machine.

Each program was run 3 times and the last time was taken (to make sure things were cached equally).


## rare

At no point scanning the data does `rare` exceed ~42MB of resident memory. Buffer sizes can be tweaked
to adjust memory usage.

```bash
$ rare -v
rare version 0.4.3, e0fc395; regex: re2

$ time rare filter -m '" (\d{3})' -e "{1}" -z testdata/*.gz | wc -l
Matched: 82,733,280 / 82,733,280
82733280

real    0m3.409s
user    0m32.750s
sys     0m2.175s
```

When aggregating data, `rare` is significantly faster than alternatives.

```bash
$ time rare histo -m '" (\d{3})' -e "{1}" -z testdata/*.gz
404         54,843,840
200         25,400,160
400         2,412,960
405         56,640
408         13,920
Matched: 82,733,280 / 82,733,280 (Groups: 8)
[8/8] 13.93 GB (4.27 GB/s)

real    0m3.283s
user    0m31.485s
sys     0m1.497s
```

And, as an alternative, using *dissect* matcher instead of regex is even slightly faster:

```bash
$ time rare histo -d '" %{CODE} ' -e '{CODE}' -z testdata/*.gz
404         54,843,840
200         25,400,160
400         2,412,960
405         56,640
408         13,920
Matched: 82,733,280 / 82,733,280 (Groups: 8)
[8/8] 13.93 GB (5.61 GB/s)

real    0m2.546s
user    0m22.922s
sys     0m1.491s
```

### pcre2

The PCRE2 version is approximately the same on a simple regular expression, but begins to shine
on more complex regex's.

```bash
# Normal re2 version
$ time rare table -z -m "\[(.+?)\].*\" (\d+)" -e "{buckettime {1} year nginx}" -e "{bucket {2} 100}" testdata/*.gz
    2020       2019
400 28,994,880 28,332,480
200 17,084,640 8,316,000
300 2,880      2,400
Matched: 82,733,280 / 82,733,280 (R: 3; C: 2)
[8/8] 13.93 GB (596.89 MB/s)

real    0m23.819s
user    3m52.252s
sys     0m1.625s

# libpcre2 version
$ time rare-pcre table -z -m "\[(.+?)\].*\" (\d+)" -e "{buckettime {1} year nginx}" -e "{bucket {2} 100}" testdata/*.gz
    2020       2019
400 28,994,880 28,332,480
200 17,084,640 8,316,000
300 2,880      2,400
Matched: 82,733,280 / 82,733,280 (R: 3; C: 2)
[8/8] 13.93 GB (2.10 GB/s)

real    0m6.813s
user    1m15.638s
sys     0m1.985s
```


## zcat & grep

```
$ time zcat testdata/*.gz | grep -Poa '" (\d{3})' | wc -l
82733280

real    0m28.414s
user    0m35.268s
sys     0m1.865s

$ time zcat testdata/*.gz | grep -Poa '" 200' > /dev/null

real    0m28.616s
user    0m27.517s
sys     0m1.658s

```

I believe the largest holdup here is the fact that zcat will pass all the data to grep via a synchronous pipe, whereas
rare can process everything in async batches.  Using `pigz` or `zgrep` instead didn't yield different results, but on single-file
results they did perform comparibly.

## Ripgrep

Ripgrep (`rg`) is the most comparible for the use-case, but lacks
the complete functionality that rare exposes.

```bash
$ time rg -z '" (\d{3})' testdata/*.gz | wc -l
82733280

real    0m7.058s
user    0m40.284s
sys     0m8.962s
```

# Other Tools

If there are other tools worth comparing, please create
a new issue on the [github tracker](https://github.com/zix99/rare/issues).
