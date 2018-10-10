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
	flag.StringVar(&e, "e", "", "function to evaluate")
	flag.Parse()
	if e == "" {
		fmt.Fprintln(os.Stderr, "Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	v, err := interp.InterpretString(e)
	if err != nil {
		errors.Explain(err, e)
		os.Exit(1)
	}
	fmt.Println(v.String())
}
