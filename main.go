package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/alecthomas/participle/lexer"
	"github.com/c-bata/go-prompt"
	"github.com/texttheater/bach/builtin"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
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
	success := executeCLI(cli.Program)
	if !success {
		os.Exit(1)
	}
}

func repl() {
	p := prompt.New(func(program string) {
		execute(program)
	}, func(prompt.Document) []prompt.Suggest {
		return nil
	},
		prompt.OptionPrefix("bach> "))
	p.Run()
}

func execute(program string) (success bool) {
	_, value, err := interpreter.InterpretString(builtin.InitialShape, states.InitialState, program)
	if err != nil {
		errors.Explain(os.Stderr, err, program)
		return false
	}
	if !printValue(value, program) {
		return false
	}
	return true
}

func executeCLI(program string) (success bool) {
	initialShape := builtin.InitialShape
	initialShape.Type = types.NewArr(types.Str{})
	initialState := states.InitialState
	initialState.Value = states.ReaderValue{Reader: os.Stdin}
	val, err := builtin.Lines(initialState, nil, nil, lexer.Position{}).Eval()
	if err != nil {
		errors.Explain(os.Stderr, err, program)
		return false
	}
	initialState.Value = val
	typ, value, err := interpreter.InterpretString(initialShape, initialState, program)
	if err != nil {
		errors.Explain(os.Stderr, err, program)
		return false
	}
	if types.AnyArr.Subsumes(typ) {
		iter := states.IterFromValue(value)
		for {
			el, ok, err := iter()
			if err != nil {
				errors.Explain(os.Stderr, err, program)
				return false
			}
			if !ok {
				break
			}
			if !printValue(el, program) {
				return false
			}
		}
	} else if (types.Null{}).Subsumes(typ) {
		// do nothing
	} else {
		if !printValue(value, program) {
			return false
		}
	}
	return true
}

func printValue(value states.Value, program string) (success bool) {
	str, err := value.Repr()
	if err != nil {
		errors.Explain(os.Stderr, err, program)
		return false
	}
	fmt.Println(str)
	return true
}
