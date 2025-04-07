# Funcer Definition Expressions

To define a funcer, use a *funcer definition expression*. It has the form `for
I def N(P1, P2, ..., Pn) O as B ok`, where `I` is the funcer’s *input type*,
`I` is its *name*, `P1, P2, ..., Pn` are `n` or more *parameters* (if `n=0`,
leave out the parentheses), `O` is the funcer’s *output type*, and `B` is its
*body*. When called, the funcer will evaluate like `B`, with parameters bound
to the arguments passed at call time.

Here are some examples of defining and then calling funcers without parameters:

```bachdoc
P for Num def plusOne Num as +1 ok 1 plusOne
T Num
V 2

P for Num def plusOne Num as +1 ok 1 plusOne plusOne
T Num
V 3
```


## Parameters

In addition to the input, funcers can take *arguments*, which are functions
that can be called in the funcer body to use their outputs. *Parameters*
specify what kind of arguments are expected. In the simplest case, a parameter
is simply a name and a type. This means that inside the funcer body, the given
function can be called with any input, by the given name, and will return a
value of the given type.

```bachdoc
P for Num def plus(b Num) Num as +b ok 1 plus(1)
T Num
V 2

P for Any def f(a Num, b Num) Arr<Num, Num> as [a, b] ok f(2, 3)
T Arr<Num, Num>
V [2, 3]
```


## Parameters with Parameters

In more complex cases, parameters can restrict arguments to be partially
applied funcer calls. In such cases, the syntax for a parameter is `for I N(P1,
P2, ..., Pn) O`. This means that inside the funcer body, for inputs of type
`I`, the given funcer will be available to call by the name `N` with arguments
as specified by the parameters `P1, P2, ..., Pn` (as before, drop the
parentheses if `n=0`) and will return a value of type `O`. Here, `P1, P2, ...,
Pn` do not contain names.

```bachdoc
P for Num def apply(for Num f Num) Num as f ok 1 apply(+1)
T Num
V 2

P for Num def apply(for Num f Num) Num as f ok 2 =n apply(+n)
T Num
V 4

P for Num def connectSelf(for Num f(for Any Num) Num) Num as =x f(x) ok 1 connectSelf(+)
T Num
V 2

P for Num def connectSelf(for Num f(for Any Num) Num) Num as =x f(x) ok 3 connectSelf(*)
T Num
V 9

P for Num def connectSelf(for Num f(Num) Num) Num as =x f(x) ok 1 connectSelf(+)
T Num
V 2
```


## Type Variables

Funcer definition expressions can use type variables, written between angle
brackets. When applying a funcer to an input expression, type variables are
bound from left to right.

```bachdoc
P for <A> def apply(for <A> f <B>) <B> as f ok 1 apply(+1)
T Num
V 2

P for <A>|Null def myMust <A> as is Null then fatal else id ok ok 123 myMust
T Num
V 123

P for <A>|Null def myMust <A> as is Null then fatal else id ok ok "abc" myMust
T Str
V "abc"

P for <A>|Null def myMust <A> as is Null then fatal else id ok ok null myMust
T <A>
E {"Kind": "Value", "Code": "UnexpectedValue", "GotValue": "null"}

P for <A>|Null def myMust <A> as is <A> then id else fatal ok ok null myMust
T Null
V null

P for <A Obj<a: Num, Any>> def f <A> as id ok {a: 1} f
T Obj<a: Num, Void>
V {a: 1}
```


## Type Variables with Bounds

Type variables can be given upper bounds, they are then constrained to match
subtypes of the given upper bound type.

```bachdoc
P for Any def f(for Any g <A Arr<Any...>>) <A> as g ok f([1, "a"])
T Arr<Num, Str>
V [1, "a"]

P for Any def f(for Any g <A Arr<Any...>>) <A> as g ok f("a")
E {"Kind": "Type", "Code": "ArgHasWrongOutputType", "ArgNum": 1, "WantType": "<A Arr<Any...>>", "GotType": "Str"}
```

## Matching Calls to Definitions

Note that funcers are looked up by input type, name, and number of paramters.
They cannot be overloaded with respect to the parameters themselves. Thus,
calling a funcer on the wrong input or with the wrong number of arguments
results in a `NoSuchFunction` error. Calling it with the wrong kinds of
arguments, by contrast, leads to an `ArgHasWrongOutputType` error.

```bachdoc
P for <A Obj<a: Num>> def f <A> as id ok {} f
E {"Kind": "Type", "Code": "NoSuchFunction", "InputType": "Obj<Void>", "Name": "f", "NumParams": 0}

P for Str def f Obj<> as {} ok "abc" reFindAll(f)
E {"Kind": "Type", "Code": "ArgHasWrongOutputType", "ArgNum": 1, "WantType": "<A Null|Obj<start: Num, 0: Str, Any>>", "GotType": "Obj<Any>"}
```
