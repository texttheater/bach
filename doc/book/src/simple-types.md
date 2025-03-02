# Simple Types

The simple types in Bach are `Null`, `Bool`, `Num`, and `Str`. The following
table shows some examples of programs that evaluate to a value of each type.

```bachdoc
P null
T Null
V null
E null

P false
T Bool
V false
E null

P true
T Bool
V true
E null

P 42
T Num
V 42
E null

P 0.3
T Num
V 0.3
E null

P "Hello world!"
T Str
V "Hello world!"
E null
```

The `Null` type has a single value, `null`.

The are two `Bool` values: `true` and `false`.

`Num` values are [IEEE754 double-precision floating-point numbers](https://en.wikipedia.org/wiki/Double-precision_floating-point_format).

`Str` values are sequences of bytes.
