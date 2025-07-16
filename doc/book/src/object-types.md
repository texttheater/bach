# Object Types

An object is a unique mapping from strings (“properties”) to values. The most
general object type is `Obj<Any>`. A more restrictive type can be specified for
all values, e.g., `Obj<Str>`. If specific properties and their types are known,
this can be part of the type too, e.g.: `Obj<a: Num, b: Num, Str>`. The type at
the end then describes all *other* values.

```bachdoc
P {}
T Obj<Void>
V {}

P {a: 1}
T Obj<a: Num, Void>
V {a: 1}

P {a: 1, b: "c"}
T Obj<a: Num, b: Str, Void>
V {a: 1, b: "c"}

P for Any def f Obj<Num> as {a: 1, b: 2} ok f
T Obj<Num>
V {a: 1, b: 2}

P for Any def f Obj<Any> as {a: 1, b: "c"} ok f
T Obj<Any>
V {a: 1, b: "c"}

P for Any def f Obj<a: Num, b: Str, Any> as {a: 1, b: "c", d: false} ok f
T Obj<a: Num, b: Str, Any>
V {a: 1, b: "c", d: false}
```
