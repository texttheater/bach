# Composition Expressions

A composition expression consists of two other expressions, separated by
whitespace. The output of the first expression is then used as the input to the
other. By adding another expression at the end, you again create a larger
composition expression, and so on.

```bachdoc
P "abc"
T Str
V "abc"

P "abc" codePoints
T Arr<Num...>
V [97, 98, 99]

P "abc" codePoints len
T Num
V 3

P "abc" codePoints len *3
T Num
V 9
```
