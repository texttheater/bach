package main

import (
	"fmt"
	"os"

	"github.com/texttheater/bach/builtin"
	"github.com/texttheater/bach/expressions"
)

var FuncersByCategory = map[string][]expressions.Funcer{
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
		fmt.Printf("## `%s`\n\n", funcer.Name)
	}
}
