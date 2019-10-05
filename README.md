# ReBlurb

A file scan/regex extractor and realtime summarizor.

Supports various CLI-based graphing and metric formats.

# Example

## Extract status codes from nginx logs

```bash
$ rare -m '"(\w{3,4}) ([A-Za-z0-9/.]+).*" (\d{3})' -e '$3 $1' h access.log
200 GET                          160663
404 GET                          857
304 GET                          53
200 HEAD                         18
403 GET                          14
```

# Output Formats

## Histogram (histo)

## Counts (count)

## Numerical Aggregats (aggr)

# License
