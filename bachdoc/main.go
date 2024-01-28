package main

import (
	"fmt"
	"os"

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
		fmt.Printf("## %s\n\n", funcer.SignatureAsMarkdown())
		fmt.Printf("%s\n\n", funcer.Summary)
		fmt.Printf("|  | Type | Value |\n")
		fmt.Printf("|---|---|---|\n")
		fmt.Printf("| Input | %s | %s |\n", funcer.InputType, funcer.InputDescription)
		for i, param := range funcer.Params {
			fmt.Printf("| `%s` (param #%d) | %s | %s |\n", param.Name, i+1, param, param.Description)
		}
		fmt.Printf("|Output | %s | %s |\n\n", funcer.OutputType, funcer.OutputDescription)
		fmt.Printf("%s\n\n", funcer.Notes)
	}
}
