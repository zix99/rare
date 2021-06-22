# Benchmarks

Below, rare is compared to various other common and popular tools on CPU user and
real time.

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
