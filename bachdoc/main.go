package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/texttheater/bach/builtin"
	"github.com/texttheater/bach/docutil"
	"github.com/texttheater/bach/shapes"
)

var FuncersByCategory = map[string][]shapes.Funcer{
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

func main() {
	ctx := kong.Parse(&cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}

var cli struct {
	Builtin    BuiltinCmd    `cmd:"" help:"Generate documentation for the builtin funcers of a given category."`
	Preprocess PreprocessCmd `cmd:"" help:"Preprocess a mdBook .md file, formatting Bach code examples"`
}

type BuiltinCmd struct {
	Category string `arg:"" help:"Funcer category."`
}

func (b *BuiltinCmd) Run() error {
	funcers, ok := FuncersByCategory[b.Category]
	if !ok {
		return errors.New("unknown category")
	}
	for _, funcer := range funcers {
		fmt.Printf("## %s\n\n", funcer.Name)
		fmt.Printf("%s\n\n", funcer.Summary)
		fmt.Printf("| | Type | Value |\n")
		fmt.Printf("|---|---|---|\n")
		fmt.Printf("| Input | %s | %s |\n", docutil.InlineCode(funcer.InputType.String()), funcer.InputDescription)
		for i, param := range funcer.Params {
			fmt.Printf("| %s (param #%d) | %s | %s |\n", param.Name, i+1, docutil.InlineCode(param.String()), param.Description)
		}
		fmt.Printf("|Output | %s | %s |\n\n", docutil.InlineCode(funcer.OutputType.String()), funcer.OutputDescription)
		fmt.Printf("%s\n\n", funcer.Notes)
		fmt.Printf("### Examples\n\n")
		docutil.PrintExamplesTable(os.Stdout, funcer.Examples)
		fmt.Printf("\n")
	}
	return nil
}

type PreprocessCmd struct {
	// TODO
}

func (p *PreprocessCmd) Run() error {
	// TODO
	return nil
}
