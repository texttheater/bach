package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/c-bata/go-prompt"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/states"
)

func help() {
	fmt.Fprintln(os.Stderr, "Usage:")
	flag.PrintDefaults()
}

func repl() {
	p := prompt.New(func(program string) {
		execute(program, true)
	}, func(prompt.Document) []prompt.Suggest {
		return nil
	},
		prompt.OptionPrefix("bach> "))
	p.Run()
}

func execute(program string, displayResult bool) (success bool) {
	_, value, err := interpreter.InterpretString(program)
	if err != nil {
		errors.Explain(err, program)
		return false
	}
	if displayResult {
		str, err := value.Repr()
		if err != nil {
			errors.Explain(err, program)
			return false
		}
		fmt.Println(str)
	} else {
		err := forceEvaluation(value)
		if err != nil {
			errors.Explain(err, program)
			return false
		}
	}
	return true
}

func main() {
	// parse command line
	var e string
	var o string
	var h bool
	flag.StringVar(&e, "e", "", "function to evaluate")
	flag.StringVar(&o, "o", "", "function to evaluate, output result")
	flag.BoolVar(&h, "h", false, "print help message and exit")
	flag.Parse()
	// help
	if h {
		help()
		os.Exit(1)
	}
	// REPL
	if e == "" && o == "" {
		repl()
		os.Exit(0)
	}
	// execute program given on command line
	var program string
	var displayResult bool
	if e != "" {
		program = e
		displayResult = false
	}
	if o != "" {
		program = o
		displayResult = true
	}
	success := execute(program, displayResult)
	if !success {
		os.Exit(1)
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
			forceEvaluation(res.Value)
		}
	}
	return nil
}
