package main

import (
	"fmt"
	"html"
	"os"
	"strings"

	"github.com/texttheater/bach/builtin"
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

func paramToMarkdown(s string) string {
	return fmt.Sprintf("<code>%s</code>", strings.ReplaceAll(html.EscapeString(s), "|", "&#124;"))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "USAGE: bachdoc CATEGORY")
		os.Exit(1)
	}
	cat := os.Args[1]
	funcers, ok := FuncersByCategory[cat]
	if !ok {
		fmt.Fprintln(os.Stderr, "ERROR: Unknown category:", cat)
		os.Exit(1)
	}
	for _, funcer := range funcers {
		fmt.Printf("## %s\n\n", funcer.Name)
		fmt.Printf("%s\n\n", funcer.Summary)
		fmt.Printf("| | Type | Value |\n")
		fmt.Printf("|---|---|---|\n")
		fmt.Printf("| Input | %s | %s |\n", paramToMarkdown(funcer.InputType.String()), funcer.InputDescription)
		for i, param := range funcer.Params {
			fmt.Printf("| %s (param #%d) | %s | %s |\n", param.Name, i+1, paramToMarkdown(param.String()), param.Description)
		}
		fmt.Printf("|Output | %s | %s |\n\n", paramToMarkdown(funcer.OutputType.String()), funcer.OutputDescription)
		fmt.Printf("%s\n\n", funcer.Notes)
		fmt.Printf("### Examples\n\n")
		fmt.Printf("| Program | Type | Value | Error |\n")
		fmt.Printf("|---|---|---|---|\n")
		for _, example := range funcer.Examples {
			var typ, val, err string
			if example.OutputType == "" {
				typ = ""
			} else {
				typ = paramToMarkdown(example.OutputType)
			}
			if example.OutputValue == "" {
				val = ""
			} else {
				val = fmt.Sprintf("`%s`", example.OutputValue)
			}
			if example.Error == nil {
				err = ""
			} else {
				err = fmt.Sprintf("```\n%s\n```", example.Error)
			}
			fmt.Printf("| `%s` | %s | %s | %s |\n", example.Program, typ, val, err)
		}
		fmt.Printf("\n")
	}
}
