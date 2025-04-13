# Math

Math expressions are evaluated using the `{! expr}` helper. They
have the following rules:

- Variables are either *non-numeric* values (eg. `abc`) or keys surrounded by brackets (eg. `[x]` or `[0]`)
- Operations (eg. `+` `-` etc), found below, follow common order-of-operations
- Parenthesis can be used to group evaluations
- The result will use the minimum number of decimals to represent the value
- The formula can't reference other expression helpers

## Examples

!!! warning
    A common error is to reference matches using the expression syntax of `{0}`, which won't
    compile. Assure you are using the formula syntax of `[x]` or `[0]`

If `x=4`, then the following will evaluate:

```
{! 2+2} => 4
{! 2 * x} => 8
{! [x] * 4} => 16
{! abs(-4)} => 4
{! (2+2)*3} => 12
{! 2(1+1) } => 4
```

As a brief example


## Operations

### Binary

| Type  | Operators   |
|-------|-------------|
| Basic | `+ - * / ^` |

**Basic Operations:** `+ - * / ^`

**Bit Shift:** `<< >>`

**Binary Operators:** `= <= >= > <`

**Binary Combiners:** `&& ||`

### Unary

For single-character unary expressions, they can be applied directly
prior to the value, eg `-x` or `!x`. For more complex expressions,
they need to be followed by a group, eg. `cos(x)` or `abs(x+2)`

- `-` Negative
- `!` Binary not
- `abs` Absolute value
- `cos`
- `sin`
- `tan`

### Formats
