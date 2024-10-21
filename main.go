package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/c-bata/go-prompt"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/states"
)

var cli struct {
	Program string `arg:"" optional:"" help:"Program to execute. If not provided, Bach's REPL will be started."`

	Quiet bool `short:"q" help:"Do not print the output value of the program."`
}

func main() {
	kong.Parse(
		&cli,
		kong.Name("bach"),
		kong.Description("An interpreter for the Bach programming language."),
	)
	// REPL
	if cli.Program == "" {
		repl()
		os.Exit(0)
	}
	// execute program given on command line
	success := execute(cli.Program, !cli.Quiet)
	if !success {
		os.Exit(1)
	}
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

func forceEvaluation(v states.Value) error {
	switch v := v.(type) {
	case *states.ArrValue:
		for v != nil {
			err := forceEvaluation(v.Head)
			if err != nil {
				return err
			}
			v, err = v.Tail.EvalArr()
			if err != nil {
				return err
			}
		}
	case states.ObjValue:
		for _, w := range v {
			val, err := w.Eval()
			if err != nil {
				return err
			}
			forceEvaluation(val)
		}
	}
	return nil
}
