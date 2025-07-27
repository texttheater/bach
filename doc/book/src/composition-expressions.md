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

Note that when you compose with expressions that ignore their input, such as
literals, only the last one has any effect on the output.

```bachdoc
P 1
T Num
V 1

P 1 2
T Num
V 2

P 1 2 3.5
T Num
V 3.5
```
