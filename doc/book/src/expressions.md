# Expressions

In Bach, everything is a function. For example:

* `2` is a function that takes any input, ignores it, and returns the number 2.
* `join(",")` is a function that takes a list of strings and concatenates them,
  using the comma as a separator.
* `+(2)` is a function that takes a number as input, adds 2 to it, and returns
  the result.
* `+2` is the same, with syntactic sugar applied.
* `if %3 ==0 then "multiple of 3" else "not a multiple of 3" ok` is a function
  that takes a number and returns different strings depending on whether the
  number is divisible by 3.
* `+2 *3` is a composition of two functions. It is also itself a function (one
  that takes a number, adds 2 to it, then multiplies the result by 3).
