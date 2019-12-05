# Why Bach?

Bach is a general-purpose programming language designed in particular for
writing “one-liners” tor text processing on the command line. It is designed to
have the following properties:

* **Concise.** One-liners must be quick to type, there is not much boiler
  plate.
* **Readable.** Bach programs consist mostly of English words. Bach makes
  minimal use of special characters. You can read a Bach program and get an
  idea of what it does even without being familiar with the language.
* **Flowing.** Bach programs are evaluated from left to right, and can be read
  and understood thus.
* **Orthogonal.** Bach’s language elements are carefully designed so that there
  is one – and preferably only one – obvious way to do any given task. This
  makes life easier for both writers and readers of Bach programs. When in
  doubt, orthogonality trumps conciseness.

Bach is designed as follows to achieve these desiderata: a Bach program is a
series of *functions*. They are just written next to each other with spaces in
between. Each function’s output is the next function’s input. In technical
terms, a Bach program is a *composition* of its component functions.
Schematically, this Bach program

    f g h

corresponds to this expression in a traditional applicative programming
language,

    h(g(f))

and to this Unix pipeline:

    f | g | h

In Bach, the concept of function composition is so central that there is no
operator for it (like `|` in Unix shells or `.` in Haskell). It’s just what you
get when you put two or more functions next to each other.

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
