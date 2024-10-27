# Object Types

An object is a unique mapping from strings to values. The most general object
type is `Obj<Any>`. A more restrictive type can be specified for all values,
e.g., `Obj<Str>`. If specific keys and their types are known, this can be part
of the type too, e.g.: `Obj<a: Num, b: Num, Str>`. The type at the end then
describes all *other* values.

{{#include ../bachdoc/examples/object-types.md}}
