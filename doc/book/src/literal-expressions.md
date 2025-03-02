# Literal Expressions

## `Num` Literals

Nonnegative real numbers can be written as integers, with a decimal point, or
in exponential notation. Negative, infinity, and NaN values can be created
using the builtin `-`, `inf`, and `nan` funcers.

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
