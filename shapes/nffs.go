package shapes

import (
	"fmt"
	"strings"
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
)

func ResolveNFF(Pos lexer.Position, inputShape Shape, name string, argFunctions []Function) (Function, error) {
	argShapes := make([]Shape, len(argFunctions))
	for i, f := range argFunctions {
		argShapes[i] = f.OutputShape(InitialShape)
	}
	stack := inputShape.Stack
Entries:
	for stack != nil {
		if !(stack.Head.Name == name) {
			stack = stack.Tail
			continue
		}
		if len(stack.Head.ArgTypes) != len(argFunctions) {
			stack = stack.Tail
			continue
		}
		if !stack.Head.InputType.Subsumes(inputShape.Type) {
			stack = stack.Tail
			continue
		}
		for i, argType := range stack.Head.ArgTypes {
			if !argType.Subsumes(argShapes[i].Type) {
				stack = stack.Tail
				continue Entries
			}
		}
		return stack.Head.Funcer(argFunctions), nil
	}
	return nil, errors.E("type", Pos, "no function found: for %v %v(%s)", inputShape.Type, name, formatArgTypes(argShapes))
}

func formatArgTypes(argShapes []Shape) string {
	formatted := make([]string, len(argShapes))
	for i, s := range argShapes {
		formatted[i] = fmt.Sprintf("%v", s.Type)
	}
	return strings.Join(formatted, ", ")
}
