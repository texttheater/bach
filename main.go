package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/states"
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
		fmt.Println(str)
	} else {
		err := forceEvaluation(value)
		if err != nil {
			errors.Explain(err, program)
			os.Exit(1)
		}
	}
}

func forceEvaluation(v states.Value) error {
	switch v := v.(type) {
	case *states.ArrValue:
		for v != nil {
			err := forceEvaluation(v.Head)
			if err != nil {
				return err
			}
			v, err = v.GetTail()
			if err != nil {
				return err
			}
		}
	case states.ObjValue:
		for _, w := range v {
			res := w.Eval()
			if res.Error != nil {
				return res.Error
			}
			forceEvaluation(res.State.Value)
		}
	}
	return nil
}
