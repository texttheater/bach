package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interpreter"
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
	_, value, err := interpreter.InterpretString(program)
	if err != nil {
		errors.Explain(err, program)
		os.Exit(1)
	}
	if o != "" {
		str, err := value.String()
		if err != nil {
			errors.Explain(err, program)
			os.Exit(1)
		}
		fmt.Println("here")
		fmt.Println(str)
	}
}
