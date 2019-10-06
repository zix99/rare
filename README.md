# rare

[![Build Status](https://travis-ci.org/zix99/rare.svg?branch=master)](https://travis-ci.org/zix99/rare)

A file scan/regex extractor and realtime summarizor.

Supports various CLI-based graphing and metric formats.

# Example

## Extract status codes from nginx logs

```bash
$ rare histo -m '"(\w{3,4}) ([A-Za-z0-9/.]+).*" (\d{3})' -e '{3} {1}' access.log
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

# Roadmap

 * Bucketing and ranging

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
