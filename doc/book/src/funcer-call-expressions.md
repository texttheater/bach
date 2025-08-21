# Funcer Call Expressions

Every Bach expression denotes a function that maps the *input* to the *output*.
*Funcers* are generalizations of functions: in addition to the input, they can
have *arguments* on which their output depends. But there are also 0-argument
funcers. A number of funcers are built-in to Bach (see [Builtin Funcer
Reference](./builtin-funcer-reference.md)), and you can define your own (see
[Funcer Definition Expressions](./funcer-definition-expressions.md)). Arguments
go between parentheses. No space is allowed between the funcer name and the
opening parenthesis. For example, `len` returns the length of an array, `join`
joins a list of strings with a given glue string, and `+(1)` adds 1 to a
number.

```bachdoc
P [1, 2, 3] len
T Num
V 3
E null

P ["a", "b", "c"] join("+")
T Str
V "a+b+c"
E null

P 1 +(1)
T Num
V 2
E null
```

## Syntactic Sugar

To make arithmetic expressions etc. less noisy, when calling the funcers `+`,
`-`, `*`, `/`, `%`, `<`, `>`, `==`, `<=`, `>=`, `**` with one argument that is
a single literal or funcer call, the parentheses can be omitted. Note that
there are no operator precedence rules in Bach – evaluation is from left to
right. To change precedence, put arguments in parentheses by not omitting the
parentheses.

```bachdoc
P 1 +1
T Num
V 2
E null

P 2 +3 *5
T Num
V 25
E null

P 2 +(3 *5)
T Num
V 17
E null
```

With a single argument that is a string, regexp, array, or object literal,
parentheses may also be omitted.

```bachdoc
P ["a", "b", "c"] join"+"
T Str
V "a+b+c"
E null
```

## Overloading

Funcers are selected according to the input type, the name, and the number of
arguments. For example, there are two builtin 1-arg funcers named `+`: one for
`Num` inputs (which expects the first argument to evaluate to a `Num` and
returns the sum) and one for `Str` inputs (which expectes the first argument to
evaluate to a `Str` and returns the concatenation).

```bachdoc
P 1 +2
T Num
V 3
E null

P "a" +"b"
T Str
V "ab"
E null

P "a" +2
E {"Kind": "Type", "Code": "ArgHasWrongOutputType", "ArgNum": 1, "WantType": "Str", "GotType": "Num"}
```

## Partial Application

In most cases, funcers are defined so that their argument can be any type of
expression as long as it has the right input and output type. But it is also
possible for funcer parameters to require a partially applied funcer call with
`n` remaining arguments to be filled in by the funcer. To call a funcer with
`m+n` arguments to fill a parameter that requires `n` open arguments, pass only
the first `m` arguments. As an example with `m=0` and `n=1`, the builtin funcer
`sort` takes a comparison funcer with an open argument as its sole argument.
Here, we use the number comparison funcer `>`:

```bachdoc
P [7, 3, 2, 5] sort(>)
T Arr<Num...>
V [7, 5, 3, 2]
E null
```
