# Benchmarks

Below, rare is compared to various other common and popular tools.

It's worth noting that in many of these results rare is just as fast, but part
of that reason is that it consumes CPU in a more efficient way (go is great at parallelization).
So take that into account, for better or worse.

All tests were done on ~83MB of gzip'd (1.5GB gunzip'd) nginx logs spread across 10 files.  They
were run on a spinning disk on an older machine. New machines run significantly faster.

Each program was run 3 times and the last time was taken (to make sure things were cached equally).


## rare

At no point scanning the data does `rare` exceed ~4MB of resident memory.

```bash
$ rare -v
rare version 0.4.3, e0fc395; regex: re2

$ time rare filter -m '" (\d{3})' -e "{1}" -z testdata/*.gz | wc -l
Matched: 8,373,328 / 8,373,328
8373328

real    0m3.266s
user    0m10.607s
sys     0m0.769s
```

When aggregating data, `rare` is significantly faster than alternatives.

```bash
$ time rare histo -m '" (\d{3})' -e "{1}" -z testdata/*.gz
404                 5,557,374 
200                 2,564,984 
400                 243,282   
405                 5,708     
408                 1,397     
Matched: 8,373,328 / 8,373,328 (Groups: 8)
[9/9] 1.41 GB (514.25 MB/s)

real    0m2.870s
user    0m9.606s
sys     0m0.393s
```

And, as an alternative, using *dissect* matcher instead of regex is even slightly faster:

```bash
$ time rare histo -d '" %{CODE} ' -e '{CODE}' -z testdata/*.gz
404         5,557,374 
200         2,564,984 
400         243,282   
405         5,708     
408         1,397     
Matched: 8,373,328 / 8,373,328 (Groups: 8)
[9/9] 1.41 GB (531.11 MB/s)

real    0m2.533s
user    0m7.976s
sys     0m0.350s
```

### pcre2

The PCRE2 version is approximately the same on a simple regular expression, but begins to shine
on more complex regex's.

```bash
# Normal re2 version
$ time rare table -z -m "\[(.+?)\].*\" (\d+)" -e "{buckettime {1} year nginx}" -e "{bucket {2} 100}" testdata/*.gz
          2020      2019      
400       2,915,487 2,892,274           
200       1,716,107 848,925             
300       290       245                 
Matched: 8,373,328 / 8,373,328 (R: 3; C: 2)
[9/9] 1.41 GB (52.81 MB/s)

real    0m27.880s
user    1m28.782s
sys     0m0.824s

# libpcre2 version
$ time rare-pcre table -z -m "\[(.+?)\].*\" (\d+)" -e "{buckettime {1} year nginx}" -e "{bucket {2} 100}" testdata/*.gz
          2020      2019      
400       2,915,487 2,892,274           
200       1,716,107 848,925             
300       290       245                 
Matched: 8,373,328 / 8,373,328 (R: 3; C: 2)
[9/9] 1.41 GB (241.82 MB/s)

real    0m5.751s
user    0m20.173s
sys     0m0.461s
```


## zcat & grep

```
$ time zcat testdata/*.gz | grep -Poa '" (\d{3})' | wc -l
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
rare can process everything in async batches.  Using `pigz` or `zgrep` instead didn't yield different results, but on single-file
results they did perform comparibly.

## Ripgrep

Ripgrep (`rg`) is the most comparible for the use-case, but lacks
the complete functionality that rare exposes.

```bash
$ time rg -z '" (\d{3})' testdata/*.gz | wc -l
8373328

real    0m3.791s
user    0m8.149s
sys     0m4.420s
```

# Other Tools

If there are other tools worth comparing, please create
a new issue on the [github tracker](https://github.com/zix99/rare/issues).
