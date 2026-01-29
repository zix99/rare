---
description: Dissect expression syntax
order: 6
depth: 1
---
# Dissect Syntax

*Dissect* is a simple token-based search algorithm, and can
be up to 10x faster than regex (and 40% faster than PCRE).

It works by searching for for constant delimiters in a string
and extracting the text between the tokens as named keys.

*rare* implements a subset of the full dissect algorithm.

**Syntax Example:**
```
prefix %{name} : %{value} - %{?ignored}
```

## Syntax

- Anything in a `%{}` is a variable token.
- A blank token, or a token that starts with `?` is skipped. eg `%{}` or `%{?skipped}`
- Tokens are extracted by both name and index (in the order they appear).
- Index `{0}` is the full match, including the delimiters
- Patterns don't need to match the entire line

## Examples

### Simple

```
prefix %{name} : %{value}
```

Will match:
```
prefix bob : 123
```

And extract 3 index-keys:
```
0: prefix bob : 123
1: bob
2: 123
```

And will extract two named keys:
```
name=bob
value=123
```

### Nginx Logs

As a simple example, to parse nginx logs that look like:

```
104.238.185.46 - - [19/Aug/2019:02:26:25 +0000] "GET / HTTP/1.1" 200 546 "-" "Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.4 (KHTML, like Gecko) Chrome/98 Safari/537.4 (StatusCake)"
```

The following dissect expression can be used:

```
%{ip} - - [%{timestamp}] "%{verb} %{path} HTTP/%{?http-version}" %{status} %{size} "-" "%{useragent}"
```

Which, as json, will return:
```json
{
    "timestamp": "12/Dec/2019:17:54:13 +0000",
    "verb": "POST",
    "path": "/temtel.php",
    "status": 404,
    "size": 571,
    "useragent": "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36",
    "ip": "203.113.174.104"
}
```
