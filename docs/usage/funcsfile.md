# Expression Functions File

A *functions file* allows you to specify additional expression
helpers that can be loaded and used in *rare*.

## Example

For example, if you frequently need to double a number you
could create a function:

```funcfile title="custom.funcs"
double {multi {0} 2}
```

And use it in rare with argument `--funcs`:
```sh
rare --funcs custom.funcs filter -m '(\d+)' -e '{double {1}}' access.log
```

Or via environment variable `RARE_FUNC_FILES`:
```sh
export RARE_FUNC_FILES=/path/to/custom.funcs
rare filter -m '(\d+)' -e '{double {1}}' access.log
```

You can load multiple files by providing `--funcs` argument multiple
times, or providing a comma-separated list to `$RARE_FUNC_FILES`

## Format

A functions file is key-value pairs of `name` to `expression`. Lines
starting with `#`, or any characters after `#`, are considered comments.

*Expressions* can be multi-line by ending the previous line with a `\`.

```funcsfile
# Allows comments that start with '#'
name-of-func {sumi {0} {1}} # comments can also go here

# Multi-line ends with '\'
classifylen {switch \
    # short string
    {lt {len {0}} 5} short \
    # long string
    {gt {len {0}} 15} long \
    medium \ # else, medium
}
```
