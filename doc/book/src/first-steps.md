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

`out` returns the special value `null`, which the Bach CLI ignores, so you you
only get that one line of output.


## Processing Lines

On the command line, it is common to process files or input streams line by
line, e.g., by applying an operation to each input line or by computing an
aggregate function over them. Bach’s CLI treats the standard input (STDIN) as a
list of lines that you can apply functions to. If those functions return a list
of strings again, those are written out line by line.

For example, you can parse each line as a number and add one to it:

    $ echo "1\n2.5\n3" | bach 'each(parseFloat +1)'
	2
	3.5
	4

Or you can sort the input lines by length;

    $ echo "aaa\na\naa" | bach 'sortBy(codePoints len, <)'
	a
	aa
	aaa


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
