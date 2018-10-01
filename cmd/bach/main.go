package main

import (
	"flag"
	"fmt"
	"os"

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
	f, err := x.Function(types.AnyType{})
	if err != nil {
		errors.Explain("type", e, err)
		os.Exit(1)
	}
	// evaluate
	output := f.Evaluate(values.NullValue{}) // TODO error handling
	fmt.Println(output.String()) // TODO sequence handling
}
