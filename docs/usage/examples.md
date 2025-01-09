# Examples

!!! note
    Please feel free to contribute your own examples on github

## Simple Text

### Histogram Values in Text File

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

### Extract Numbers from Text
```sh
$ rare filter --match "(\d+)" input.txt
```

### Extract matched value as JSON
```sh
$ rare f --match "(?P<val>\d+)" -e "{.}" simple.log
{"val": 1}
{"val": 2}
{"val": 1}
```

### Histogram of Numbers in Text
```sh
$ rare histo --match "(\d+)" -e "{1}" -x input.txt
1                   3         
0                   1         
2                   1         
3                   1         

Matched: 6 / 6 (Groups: 4)

# Or with Bars/percentages
./rare histo --match "(\d+)" -e "{1}" -x simple.log
1                   7          [58.3%] ||||||||||||||||||||||||||||||||||||||||||||||||||
3                   2          [16.7%] ||||||||||||||
2                   2          [16.7%] ||||||||||||||
0                   1          [ 8.3%] |||||||

Matched: 12 / 12 (Groups: 4)

# Logarithmic Scale
./rare histo --match "(\d+)" -e "{1}" --scale log10 -x simple.log
1                   7          [58.3%] ||||||||||||||||||||||||||||||||||||||||||
3                   2          [16.7%] |||||||||||||||
2                   2          [16.7%] |||||||||||||||
0                   1          [ 8.3%] 

Matched: 12 / 12 (Groups: 4)
```

## Nginx

Example line:

```log
123.123.123.123 - - [19/Aug/2019:02:26:25 +0000] "GET / HTTP/1.1" 200 546 "-" "Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.4 (KHTML, like Gecko) Chrome/98 Safari/537.4 (StatusCake)"
```

### Highlight / Extract HTTP Code and Size

```sh
# Will colorize HTTP code and size in full log
$ rare filter -m "(\d{3}) (\d+)" access.log

# Will only display http code and size
$ rare filter -m "(\d{3}) (\d+)" -e "{1} {2}" access.log
```

### HTTP Status Histogram

Parse error codes and graph in a histogram

```sh
$ rare h -m "\" (\d+)" -e "{1}" -z -x testdata/*

404                 5,557,374  [66.4%] ||||||||||||||||||||||||||||||||||||||||||||||||||
200                 2,564,984  [30.6%] |||||||||||||||||||||||
400                 243,282    [ 2.9%] ||
405                 5,708      [ 0.1%]
408                 1,397      [ 0.0%]
Matched: 8,373,328 / 8,383,717 (Groups: 8)
```

### Extracting Page Sizes

Page sizes, ignoring 0-sized pages

```sh
$ rare h -m "\" (\d+) (\d+)" -e "{bytesize {bucket {2} 1024}}" -i "{lt {2} 1024}" -z -x testdata/*

234 KB              3,602      [14.6%] ||||||||||||||||||||||||||||||||||||||||||||||||||
149 KB              2,107      [ 8.5%] |||||||||||||||||||||||||||||
193 KB              1,519      [ 6.2%] |||||||||||||||||||||
192 KB              1,470      [ 6.0%] ||||||||||||||||||||
191 KB              1,421      [ 5.8%] |||||||||||||||||||
Matched: 24,693 / 8,383,717 (Groups: 96) (Ignored: 8,348,635)
```

### Extract number of bytes sent by bucket, and format

This shows an example of how to bucket the values into size of `1000`. In this case, it doesn't make
sense to see the histogram by number of bytes, but we might want to know the ratio of various orders-of-magnitudes.

```sh
$ rare histo -m '"(\w{3,4}) ([A-Za-z0-9/.]+).*" (\d{3}) (\d+)' -e "{bucket {4} 10000}" -n 10 -b access.log
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

### Table of URLs to HTTP Status

Know how your URLs are responding by their http statuses

```sh
$ rare t -m "\"(\w+) (.+).*\" (\d+) (\d+)" -e "{$ {3} {substr {2} 0 20}}" -z testdata/*

                     404                  200                  400
/ HTTP/1.1           0                    127,624              5,681
/ HTTP/1.0           0                    5,222                0
/test.php HTTP/1.1   3,241                0                    0
/1.php HTTP/1.1      2,508                0                    0
/qq.php HTTP/1.1     1,908                0                    0
/index.php HTTP/1.1  1,776                0                    0
/shell.php HTTP/1.1  1,750                0                    0
/cmd.php HTTP/1.1    1,588                0                    0
/x.php HTTP/1.1      1,573                0                    0
/log.php HTTP/1.1    1,261                0                    0
/confg.php HTTP/1.1  1,253                0                    0
/api.php HTTP/1.1    1,241                0                    0
/ss.php HTTP/1.1     1,233                0                    0
/mirror/distros/vlit 0                    1,122                0
/robots.txt HTTP/1.1 1,056                0                    0
/vendor/phpunit/phpu 1,055                0                    0
/aaa.php HTTP/1.1    954                  0                    0
/hell.php HTTP/1.1   948                  0                    0
/z.php HTTP/1.1      948                  0                    0
Matched: 465,348 / 470,163 (R: 2396; C: 8)
```

### Heatmap of Requests per Week

```
$ rare heatmap -m '\[(.+?)\].*" (\d+)' -e "{timeattr {time {1}} yearweek}" -e "{2}" access.log

    - 0    4 22,602    9 45,204
    2019-34..2019-43..2019-52..2020-17..2020-25..2020-33...2020-9
200 11--111111-11111-111111111-1111111-111111-11111-1111-11--1-1-
206 -------------------------------------------------------------
301 -------------------------------------------------------------
304 -------------------------------------------------------------
400 -------------------------------------------------------------
404 33415264914143253212111-1411-13-131-1312-122111--13-1-1-13211
405 -------------------------------------------------------------
408 -------------------------------------------------------------
Matched: 1,035,666 / 1,035,666 (R: 8; C: 61)
```

### Bargraph status codes per year

**NOTE:** For stacking (`-s`), the results will be color-coded (not shown here)

```sh
$ rare bars -z -m "\[(.+?)\].*\" (\d+)" -e "{buckettime {1} year}" -e "{2}" -s testdata/*

        0 200  1 206  2 301  3 304  4 400  5 404  6 405  7 408
2019  000000000555555555555555555555555555555  3,742,444
2020  0000000000000000004455555555555555555555555555555  4,631,884
Matched: 8,374,328 / 8,384,811
```
