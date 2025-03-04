# Literal Expressions

## `Num` Literals

Nonnegative real numbers can be written in the decimal system: as integers,
with a decimal point, or in exponential notation.

```bachdoc
P 123
T Num
V 123 
E null

P 1.23
T Num
V 1.23 
E null

P 01.23
T Num
V 1.23
E null

P .23
T Num
V 0.23
E null

P 1.
T Num
V 1
E null

P 1.23e2
T Num
V 123
E null

P 123E2
T Num
V 12300
E null

P 123E+2
T Num
V 12300
E null

P 1e-1
T Num
V 0.1
E null

P .1e0
T Num
V 0.1
E null

P 0010e-2
T Num
V 0.1
E null

P 0e+5
T Num
V 0
E null

```

There are no literals for negative, infinity, or NaN values. But they can be created using the builtin `-`, `inf`, and `nan` funcers. For details, see [Math Funcers](math-funcers.md).

```bachdoc
P -1
T Num
V -1
E null

P -0010e-2
T Num
V -0.1
E null

P -0
T Num
V -0
E null

P inf
T Num
V inf
E null

P -inf
T Num
V -inf
E null

P nan
T Num
V nan
E null
```

## `Str` Literals

Strings are sequences of bytes. String literals are delimited by double quotes. Characters inside generally represent the byte sequence that is their UTF-8 encoding. For example:

* `a` represents the byte sequence `61` (hexadecimal notation; UTF-8 encoding of Latin Small Letter A)
* `~` represents the byte sequence `7E` (UTF-8 encoding of Tilde)
* `abc` represents the byte sequence `61 62 63`
* `日本語` represents the byte sequence `E6 97 A5 E6 9C AC E8 AA 9E`

There are, however, the following exceptions:

* `\a` represents the byte sequence `07` (UTF-8 encoding of Bell character)
* `\b` represents the byte sequence `08` (UTF-8 encoding of Backspace)
* `\f` represents the byte sequence `0C` (UTF-8 encoding of Form feed)
* `\n` represents the byte sequence `0A` (UTF-8 encoding of Line feed)
* `\r` represents the byte sequence `0D` (UTF-8 encoding of Carriage return)
* `\t` represents the byte sequence `09` (UTF-8 encoding of Horizontal tab)
* `\v` represents the byte sequence `09` (UTF-8 encoding of Vertical tab)
* `\\` represents the byte sequence `5C` (UTF-8 encoding of Backslash)
* `\"` represents the byte sequence `22` (UTF-8 encoding of Quotation mark)
* `\` followed by three octal digits in the range from `000` to `3ff` (inclusive) represents the corresponding byte
* `\x` followed by two hexadecimal digits represents the corresponding byte
* `\u` followed by four hexadecimal digits represents the UTF-8 encoding of the corresponding code point, if defined
* `\U` followed by eight hexadecimal digits represents the UTF-8 encoding of the corresponding code point, if defined
* `{{` represents the byte sequence `7B` (UTF-8 encoding of Left curly bracket)
* `}}` represents the byte sequence `7D` (UTF-8 encoding of Right curly bracket)
* `{`, followed by a Bach expression, followed by `}`, represents the UTF-8 encoding of what the expression evaluates to
* Other uses of `\`, `"`, `{`, or `}` inside the delimiting quotes are invalid, as is the line feed character

```bachdoc
P "a"
T Str
V "a"
E null

P "\a"
T Str
V "\a"
E null

P "\"\\a\""
T Str
V "\"\\a\""
E null

P "\141"
T Str
V "a"
E null

P "\x61"
T Str
V "a"
E null

P "\u65e5\u672c\u8a9e"
T Str
V "日本語"
E null

P "\U000065e5\U0000672c\U00008a9e"
T Str
V "日本語"
E null

P "{{}}"
T Str
V "{{}}"
E null

P "1 + 1 = {1 +1}"
T Str
V "1 + 1 = 2"
E null

P "{ {a: 1 +1} }"
T Str
V "{{a: 2}}"
E null
```

## `Arr` Literals

Array literals are delimited by square brackets. Inside, a comma-separated
sequence of Bach expressions represents the elements. An expression
representing an array can be appended as a suffix using a semicolon.

```bachdoc
P []
T Arr<>
V []
E null

P [1]
T Arr<Num>
V [1]
E null

P [1, 2, 3]
T Arr<Num, Num, Num>
V [1, 2, 3]
E null

P [1, "a"]
T Arr<Num, Str>
V [1, "a"]
E null

P [[1, 2], ["a", "b"]]
T Arr<Arr<Num, Num>, Arr<Str, Str>>
V [[1, 2], ["a", "b"]]
E null

P [1 +1]
T Arr<Num>
V [2]
E null

P [1;[]]
T Arr<Num>
V [1]
E null

P [1, 2;[3, 4]]
T Arr<Num, Num, Num, Num>
V [1, 2, 3, 4]
E null

P [3, 4] =rest [1, 2;rest]
T Arr<Num, Num, Num, Num>
V [1, 2, 3, 4]
E null

P [1, 2;[1, 2] each(+2)]
T Arr<Num, Num, Num...>
V [1, 2, 3, 4]
E null
```

## `Obj` Literals

Obj literals are delimited by curly braces. Inside, each elements consists of a
key, followed by a colon, followed by a value. Elements are separated by
commas. The order of elements does not matter. Keys are always strings, but
they can be written as identifiers or number literals; they are converted to
strings automatically.

```bachdoc
P {}
T Obj<Void>
V {}
E null

P {"a": 1}
T Obj<a: Num, Void>
V {a: 1}
E null

P {a: 1}
T Obj<a: Num, Void>
V {a: 1}
E null

P {1: "a"}
T Obj<1: Str, Void>
V {1: "a"}
E null

P {"1": "a"}
T Obj<1: Str, Void>
V {1: "a"}
E null

P {a: 1, b: "c"}
T Obj<a: Num, b: Str, Void>
V {a: 1, b: "c"}
E null

P {b: 1 +1}
T Obj<b: Num, Void>
V {b: 2}
E null
```
