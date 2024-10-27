# Array Types

The type of an array is written as `Arr<` … `>`. Between the angle brackets,
the types of the array's elements are listed. If the last element type is
followed by an ellipsis, this means that the length of the array is unknown and
that there can be 0, 1, or more elements of this type here. Array functions
such as `each(...)` usually return arrays with an unknown number of elements.

{{#include ../bachdoc/examples/array-types.md}}
