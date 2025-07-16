# Types

Bach has null, boolean, number, string, array, object, and union types.
Additionally, there is the `Void` type (there are no values of this type) and
the `Any` type (all values are of this type). If you are wondering what the
output type of your program is, you can always compose it with the `type`
function, and you will get a string representation of the output type.

The following table demonstrates this: in the Program column, there are a
number of Bach programs, each consisting of some expression followed by the
`type` function. The Type column indicates the output type of the whole
program. Here, it’s always `Str` because that is the output type of `type`.
Finally, the Value column contains the strings representing the types of the
various expressions.

```bachdoc
P null type
T Str
V "Null"

P false type
T Str
V "Bool"

P true type
T Str
V "Bool"

P 42 type
T Str
V "Num"

P 0.3 type
T Str
V "Num"

P "Hello world!" type
T Str
V "Str"

P [] type
T Str
V "Arr<>"

P {} type
T Str
V "Obj<Void>"

P for Any def f Num|Str as 0 ok f type
T Str
V "Num|Str"

P for Any def f Any as null ok f type
T Str
V "Any"
```
