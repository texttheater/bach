package main

import (
	"errors"
	"fmt"
	"html"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/texttheater/bach/builtin"
	"github.com/texttheater/bach/interpreter"
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

var ExampleSetsByName = map[string][]shapes.Example{
	"simple-types": interpreter.SimpleTypeExamples,
	"array-types":  interpreter.ArrayTypeExamples,
}

func main() {
	ctx := kong.Parse(&cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}

var cli struct {
	Builtin  BuiltinCmd  `cmd:"" help:"Generate documentation for the builtin funcers of a given category."`
	Examples ExamplesCmd `cmd:"" help:"Format a given example set as a markdown table."`
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
		fmt.Printf("| Input | %s | %s |\n", inlineCode(funcer.InputType.String()), funcer.InputDescription)
		for i, param := range funcer.Params {
			fmt.Printf("| %s (param #%d) | %s | %s |\n", param.Name, i+1, inlineCode(param.String()), param.Description)
		}
		fmt.Printf("|Output | %s | %s |\n\n", inlineCode(funcer.OutputType.String()), funcer.OutputDescription)
		fmt.Printf("%s\n\n", funcer.Notes)
		fmt.Printf("### Examples\n\n")
		printExamplesTable(funcer.Examples)
		fmt.Printf("\n")
	}
	return nil
}

type ExamplesCmd struct {
	Name string `arg:"" help:"Name of example set."`
}

func (e *ExamplesCmd) Run() error {
	examples, ok := ExampleSetsByName[e.Name]
	if !ok {
		return errors.New("unknown example set")
	}
	printExamplesTable(examples)
	return nil
}

// inlineCode takes a string representing some program code and converts it to
// a Markdown representation suitable for processing by mdbook.
func inlineCode(s string) string {
	// handle empty string specially: do not generate <code></code> tags
	if s == "" {
		return ""
	}
	// escape HTML special characters
	s = html.EscapeString(s)
	// escape characters that mdbook treats specially
	s = strings.ReplaceAll(s, "|", "&#124;")
	s = strings.ReplaceAll(s, "\\", "&#92;")
	// wrap in code tags
	s = "<code>" + s + "</code>"
	// return
	return s
}

func printExamplesTable(examples []shapes.Example) {
	fmt.Printf("| Program | Type | Value | Error |\n")
	fmt.Printf("|---|---|---|---|\n")
	for _, example := range examples {
		var err string
		if example.Error == nil {
			err = ""
		} else {
			err = strings.TrimSpace(fmt.Sprintf("%s", example.Error))
		}
		fmt.Printf(
			"| %s | %s | %s | %s |\n",
			inlineCode(example.Program),
			inlineCode(example.OutputType),
			inlineCode(example.OutputValue),
			inlineCode(err),
		)
	}
}
