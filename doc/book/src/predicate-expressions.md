# Predicate Expressions

Predicates are related to [conditional
expressions](conditional-expressions.md). A predicate is what you get when you
use an `if` clause, an `is` clause, or an `is ... with` clause on its own
(usually, as an argument to a funcer). If the condition is satisfied, it
returns `{yes: I}`, and if not, `{no: I}`, where `I` is the input. That is,
predicates are auto-completed by Bach indto conditionals with `then {yes: id}
else {no: id}`. Such predicate expressions can be passed, e.g., to some [array
funcers](./array-funcers.md).

```bachdoc
P is Null
E {"Kind": "Type", "Code": "UnreachableElseClause"}

P 2 is Num with >3
T Obj<yes: Num>|Obj<no: Num>
V {no: 2}

P 4 is Num with >3
T Obj<yes: Num>|Obj<no: Num>
V {yes: 4}

P 2 if >3
T Obj<yes: Num>|Obj<no: Num>
V {no: 2}

P 4 if >3
T Obj<yes: Num>|Obj<no: Num>
V {yes: 4}

P for Any def f Num|Str as 2 ok f is Num _
T Obj<yes: Num>|Obj<no: Str>
V {yes: 2}
```
