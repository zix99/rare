# Json

Syntax: `{json field expression}`

Extract a JSON value based on the expression statement
from [gjson](https://github.com/tidwall/gjson)

When using, you likely often want to extract a specific bit of
json from the line.  If you want to match the entire line you
will likely want to leave the match at the default `.*`, and
provide `{0}` as the field.

For example, this command would extract lastname from below:

`rare filter -e '{json {0} "name.last"}`

Example expressions from *gjson* documentation:

```json
"name.last"          >> "Anderson"
"age"                >> 37
"children"           >> ["Sara","Alex","Jack"]
"children.#"         >> 3
"children.1"         >> "Alex"
"child*.2"           >> "Jack"
"c?ildren.0"         >> "Sara"
"fav\.movie"         >> "Deer Hunter"
"friends.#.first"    >> ["Dale","Roger","Jane"]
"friends.1.last"     >> "Craig"
```

With the given example:

```json
{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
  ]
}
```
