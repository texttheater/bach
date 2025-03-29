package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/texttheater/bach/builtin"
	"github.com/texttheater/bach/shapes"
)

var funcersByCategory = map[string][]shapes.Funcer{
	"null":    builtin.NullFuncers,
	"io":      builtin.IOFuncers,
	"logic":   builtin.LogicFuncers,
	"math":    builtin.MathFuncers,
	"text":    builtin.TextFuncers,
	"arr":     builtin.ArrFuncers,
	"obj":     builtin.ObjFuncers,
	"types":   builtin.TypeFuncers,
	"values":  builtin.ValueFuncers,
	"regexp":  builtin.RegexpFuncers,
	"control": builtin.ControlFuncers,
}

type builtinCmd struct {
	Category string `arg:"" help:"funcer category"`
}

func (b *builtinCmd) Run() error {
	funcers, ok := funcersByCategory[b.Category]
	if !ok {
		return errors.New("unknown category")
	}
	for _, funcer := range funcers {
		fmt.Printf("## %s\n\n", funcer.Name)
		fmt.Printf("%s\n\n", funcer.Summary)
		fmt.Printf("| | Type | Value |\n")
		fmt.Printf("|---|---|---|\n")
		fmt.Printf("| Input | %s | %s |\n", inlineCode(funcer.InputType.String()), funcer.InputDescription)
		for i, param := range funcer.Params {
			fmt.Printf("| %s (param #%d) | %s | %s |\n", param.Name, i+1, inlineCode(param.String()), param.Description)
		}
		fmt.Printf("|Output | %s | %s |\n\n", inlineCode(funcer.OutputType.String()), funcer.OutputDescription)
		fmt.Printf("%s\n\n", funcer.Notes)
		fmt.Printf("### Examples\n\n")
		printExamplesTable(os.Stdout, funcer.Examples)
		fmt.Printf("\n")
	}
	return nil
}
