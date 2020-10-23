# Examples

Please feel free to contribute your own examples on github

## Nginx

### HTTP Status

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

### Table of URLs to HTTP Status

Know how your URLs are responding by their http statuses

```sh
$ rare t -m "\"(\w+) (.+).*\" (\d+) (\d+)" -e "{$ {3} {substr {2} 0 20}}" -z testdata/*

                     404                  200                  400                  405                  304                  301                  
/ HTTP/1.1           0                    21,614               38                   22                   8                    0                                         
/mirror/distros/vlit 0                    1,122                0                    0                    0                    0                                         
/test.php HTTP/1.1   909                  0                    0                    0                    0                    0                                         
/ HTTP/1.0           0                    833                  0                    4                    0                    0                                         
/1.php HTTP/1.1      792                  0                    0                    0                    0                    0                                         
/qq.php HTTP/1.1     611                  0                    0                    0                    0                    0                                         
../../mnt/custom/Pro 0                    0                    558                  0                    0                    0                                         
/cmd.php HTTP/1.1    511                  0                    0                    0                    0                    0                                         
/x.php HTTP/1.1      497                  0                    0                    0                    0                    0                                         
/shell.php HTTP/1.1  478                  0                    0                    0                    0                    0                                         
/log.php HTTP/1.1    399                  0                    0                    0                    0                    0                                         
/confg.php HTTP/1.1  396                  0                    0                    0                    0                    0                                         
/api.php HTTP/1.1    385                  0                    0                    0                    0                    0                                         
/ss.php HTTP/1.1     379                  0                    0                    0                    0                    0                                         
/index.php HTTP/1.1  349                  0                    0                    0                    0                    0                                         
/aaa.php HTTP/1.1    302                  0                    0                    0                    0                    0                                         
/hell.php HTTP/1.1   300                  0                    0                    0                    0                    0                                         
/z.php HTTP/1.1      300                  0                    0                    0                    0                    0                                         
/123.php HTTP/1.1    297                  0                    0                    0                    0                    0                                         
Matched: 116,729 / 117,170 (R: 1149; C: 6)
```