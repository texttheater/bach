# Conditional Expressions

Conditional expressions come in three flavors:

* `if A then B else C ok` (checks a boolean condition A)
* `is A then B else C ok` (matches the input against the pattern A)
* `is A with D then B else C ok` (matches the input against the pattern A and checks the boolean condition D)

A pattern can be a name, a type optionally combined with a name, an array
literal where the elements are patterns, or an object literal where the values
are patterns.

To check for additional alternative conditions, add alternative clauses before
the final `else` clause. Alternative clauses, too, come in three flavors. All
can be freely combined:

* `elis A then B` (checks a boolean condition A)
* `elis A then B` (matches the input against the pattern A)
* `elis A with D then B` (matches the input against the pattern A and checks the boolean condition D)

Names bound by a pattern can be used inside the following `then` clause.

```bachdoc
P if true then 2 else 3 ok
T Num
V 2

P for Num def heart Bool as if <3 then true else false ok ok 2 heart
T Bool
V true

P for Num def heart Bool as if <3 then true else false ok ok 4 heart
T Bool
V false

P for Num def expand Num as if <0 then -1 elif >0 then +1 else 0 ok ok 0 -1 expand
T Num
V -2

P for Num def expand Num as if <0 then -1 elif >0 then +1 else 0 ok ok 1 expand
T Num
V 2

P for Num def expand Num as if <0 then -1 elif >0 then +1 else 0 ok ok 0 expand
T Num
V 0
```


## Predicates

`if` clauses, `is` clauses, and `is ... with` clauses can be used on their own
(usually, as arguments to funcers). If the condition is satisfied, they return
`{yes: I}`, and if not, `{no: I}`, where `I` is the input. Such predicate
expressions can be passed, e.g., to some [array funcers](./array-funcers.md).

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
