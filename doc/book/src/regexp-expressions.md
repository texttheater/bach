# Regexp Expressions

In Bach, a *regexp* is a function that takes a `Str` and returns either `null`
or an object with at least a property `start` with type `Num`, and `0` with
type `Str`.

Regexp expressions denote such functions. They have the syntax of [re2 regular
expressions](https://github.com/google/re2/wiki/Syntax), except for `\C`,
delimited by `~`. They search the input string for a substring that matches the
regular expression. If none can be found, `null` is returned. Otherwise, an
object is returned with `start` indicating the offset (in codepoints) of the
first occurrence, and `0` containing the matched substring. If the regexp
contains capturing groups, the substrings matched by the groups are included as
additional properties – by group number and, in the case of named groups,
additionally by name.

```bachdoc
P "abccd" ~b(?P<cs>c*)d~
T Null|Obj<start: Num, 0: Str, 1: Null|Str, cs: Null|Str, Void>
V {start: 1, 0: "bccd", 1: "cc", cs: "cc"}

P "def" ~^b(?P<cs>c*)d~
T Null|Obj<start: Num, 0: Str, 1: Null|Str, cs: Null|Str, Void>
V null

P "abccd" ~^b(?P<cs>*)d~
E {"Kind": "Syntax", "Code": "BadRegexp"}
```

Regexp expressions can be used with [regexp funcers](regexp-funcers.md) for
greater effect.
