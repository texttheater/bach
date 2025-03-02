# Array Types

The type of an array is written as `Arr<` … `>`. Between the angle brackets,
the types of the array's elements are listed. If the last element type is
followed by an ellipsis, this means that the length of the array is unknown and
that there can be 0, 1, or more elements of this type here. Array functions
such as `each(...)` usually return arrays with an unknown number of elements.

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

P [1, 2, 3] each(+1)
T Arr<Num...>
V [2, 3, 4]
E null

P [1, "a"]
T Arr<Num, Str>
V [1, "a"]
E null

P [[1, 2], ["a", "b"]]
T Arr<Arr<Num, Num>, Arr<Str, Str>>
V [[1, 2], ["a", "b"]]
E null

P [1;[]]
T Arr<Num>
V [1]
E null

P [1, 2; [3, 4]]
T Arr<Num, Num, Num, Num>
V [1, 2, 3, 4]
E null

P [3, 4] =rest [1, 2; rest]
T Arr<Num, Num, Num, Num>
V [1, 2, 3, 4]
E null

P [1, 2; [1, 2] each(+2)]
T Arr<Num, Num, Num...>
V [1, 2, 3, 4]
E null

P for Arr<Any...> def f Arr<Any...> as id ok [] f
T Arr<Any...>
V []
E null
```
