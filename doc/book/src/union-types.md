# Union Types

A union type indicates that a value could be of any of two or more types.

```bachdoc
P [1] +["a"]
T Arr<Num|Str...>
V [1, "a"]
E null

P [1] +["a"] get(0)
T Num|Str
V 1
E null
```
