# Assignment Expressions

An assignment expressions consists of the character `=` immediately followed by
a pattern. It binds the names in the pattern to the corresponding parts of its
input value so you can reuse them later in the program. A pattern can be a
name, an array literal where the elements are patterns, or an object literal
where the values are patterns.

```bachdoc
P 1 +1 =a 3 *2 +a
T Num
V 8

P 1 +1 ==2 =p 1 +1 ==1 =q p ==q not
T Bool
V true

P [1, 2, 3] =[a, b, c] a
T Num
V 1

P [1, 2, 3] =[a;r] r
T Arr<Num, Num>
V [2, 3]

P {a: 1, b: 2, c: 3} ={a: d, b: e, c: f} d
T Num
V 1

P {a: 1, b: 2, c: 3} ={a: d, b: e, c: f} e
T Num
V 2

P {a: 1, b: 2, c: 3} ={a: d, b: e, c: f} f
T Num
V 3

P for Num def cube Num as =n *n *n ok 3 cube
T Num
V 27
```

An “impossible match” error occurs if the pattern cannot match any values of
the input type, e.g., if an array pattern has a different length from the
input, or matching objects with array patterns, or vice versa, or if an object
pattern contains keys that the input doesn’t.

```bachdoc
P [1, 2, 3] =[a, b]
E {"Kind": "Type", "Code": "ImpossibleMatch"}

P {a: 1, b: 2, c: 3} =[a, b]
E {"Kind": "Type", "Code": "ImpossibleMatch"}

P {a: 1, b: 2, c: 3} ={g: h}
E {"Kind": "Type", "Code": "ImpossibleMatch"}
```

A “nonexhaustive match” error occurs if the pattern can match some but not all
values of the input type, e.g., if a variable-length array type is matched as
fixed-length, or if an object type is matched against a key it might or might
not contain.

```bachdoc
P for Arr<Num...> def f Num as =[a, b] a ok
E {"Kind": "Type", "Code": "NonExhaustiveMatch"}

P for Obj<a: Num, Num> def f Num as ={a: a, b: b} a +b ok
E {"Kind": "Type", "Code": "NonExhaustiveMatch"}

P for Obj<a: Num, Num> def f Num as ={b: b} b ok
E {"Kind": "Type", "Code": "NonExhaustiveMatch"}
```
