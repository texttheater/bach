# Getter Expressions

A getter expressions consists of the `@` character followed by an
[object](object-types.md) property name or an [array](array-types.md) index. It
retrieves the value of the property or the element at the index for an input
object or array. The input type must guarantee the existence of the
property/index, otherwise a type error is raised. If the input type cannot
guarantee the existence of the property/index, use the `get` funcer [for
objects](object-funcers.md#get) or [for arrays](array-funcers.md#get).

```bachdoc
P {a: 1, b: 2} @a
T Num
V 1

P {a: 1, b: 2} @b
T Num
V 2

P {a: 1, b: 2} @c
E {"Kind": "Type", "Code": "NoSuchProperty", "WantType": "Obj<c: Any>", "GotType": "Obj<a: Num, b: Num, Void>"}

P {0: "a"} @0
T Str
V "a"

P ["a", "b", "c"] @0
T Str
V "a"

P ["a", "b", "c"] @1
T Str
V "b"

P ["a", "b", "c"] @2
T Str
V "c"

P ["a", "b", "c"] @3
E {"Kind": "Type", "Code": "NoSuchIndex"}

P ["a", "b", "c"] @-1
E {"Kind": "Type", "Code": "BadIndex"}

P ["a", "b", "c"] @1.5
E {"Kind": "Type", "Code": "BadIndex"}

P "abc" @1
E {"Kind": "Type", "Code": "NoGetterAllowed"}

P 24 @1
E {"Kind": "Type", "Code": "NoGetterAllowed"}

P for Any def f Arr<Any...> as [] ok f @1
E {"Kind": "Type", "Code": "NoGetterAllowed"}
```
