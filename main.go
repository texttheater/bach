package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interp"
)

func main() {
	var e string
	var o string
	flag.StringVar(&e, "e", "", "function to evaluate")
	flag.StringVar(&o, "o", "", "function to evaluate, output result")
	flag.Parse()
	if (e == "") == (o == "") { // exactly one must be given
		fmt.Fprintln(os.Stderr, "Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	var program string
	if e != "" {
		program = e
	}
	if o != "" {
		program = o
	}
	_, value, err := interp.InterpretString(program)
	if err != nil {
		errors.Explain(err, program)
		os.Exit(1)
	}
	if o != "" {
		fmt.Println(value)
	}
}
