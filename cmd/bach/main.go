package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/texttheater/bach/contexts"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/grammar"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func main() {
	var e string
	flag.StringVar(&e, "e", "", "function to evaluate")
	flag.Parse()
	if e == "" {
		fmt.Fprintln(os.Stderr, "Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	// parse
	x, err := grammar.Parse(e)
	if err != nil {
		errors.Explain("syntax", e, err)
		os.Exit(1)
	}
	// type-check
	c := contexts.Context{types.AnyType{}}
	f := x.Function(c) // TODO error handling
	// evaluate
	output := f.Evaluate(values.NullValue{}) // TODO error handling
	fmt.Println(output.String()) // TODO sequence handling
}
