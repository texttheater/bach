# First Steps

## The CLI

Call Bach and pass a Bach program to it as an argument on the command line. In
our first example, our program consists of the string `"Hello world!"`. Bach
string literals use double quotes, so surround your program with single quotes.

    $ bach '"Hello world!"'
    "Hello world!"

Our program creates the string `"Hello world!"`, and Bach shows this result to
us, formatted as a string literal. If we want to print the message without the
quotes, we can compose our pogram with the function `out` to print it out:

    $ bach '"Hello world!" out'
    Hello world!
    "Hello world!"

`out` returns its input value unchanged, so we are still shown the string with the quotes at the end.
To suppress this, we can call Bach with the `-q` flag:

    $ bach -q '"Hello world!" out'
    Hello world!

## The REPL

Let us now use Bach in interactive mode by using its read-eval-print loop
(REPL). To start it, we call Bach without an argument. Now we are shown the
Bach prompt:

    $ bach
    bach>

Let us again create the string `"Hello world!"`:

    bach> "Hello world!"
    "Hello world!"

Now let us compose this program with the function `codePoints`, which gives us
the list of Unicode code points in the string:

    bach> "Hello world!" codePoints
    [72, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 33]

Let us add another function, `len`, to compute the length of the string (in
terms of code points):

    bach> "Hello world!" codePoints len
    12

As a final example in this introduction, let us add 1 and 2:

    bach> 1 +(2)
    3

Bach does not have infix operators. `+` is just a function that takes an input
value (`1`) and an argument (`2`). Arguments are given in parentheses. However,
for the mathematical operators, you can leave the parentheses out. Note
however, that you can't have a space between the `+` and the `2`:

    bach> 1 +2
    3

To exit the REPL, press Ctrl+D.
