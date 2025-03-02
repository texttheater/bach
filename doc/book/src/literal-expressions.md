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

String literals are delimited by double quotes. Characters inside represent the
byte sequence that is their UTF-8 encoding, with the exception of the four
characters `\"{}`, which can be written as `\\`, `\"`, `{{`, and `}}`,
respectively. Escape sequences of the form TODO represent the UTF-8 encoding of
the Unicode codepoint TODO. Escape sequences of the form TODO represent the
literal byte TODO. Bach expressions between curly braces represent the UTF-8
encoding of the default string representation of their return value.

## `Arr` Literals

Array literals are delimited by square brackets. Inside, a comma-separated
sequence of Bach expressions represents the elements. An expression
representing an array can be appended as a suffix using a semicolon.

## `Obj` Literals

Obj literals are delimited by curly braces. Inside, each elements consists of a
key, followed by a colon, followed by a value. Elements are separated by
commas. The order of elements does not matter. Keys are always strings, but
they can be written as identifiers or number literals; they are converted to
strings automatically.
