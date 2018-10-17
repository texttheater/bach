/*
Package shapes implements shapes.

A shape consists of a type and a stack of available NFFs.

Interpreting a Bach program involves assigning each expression an input shape,
a function and an output shape. The first expression in the program gets the
initial shape, consisting of the Any type and a stack consisting only of
builtin NFFs. The input shape of an expression and the expression together
determine its function. The function and the input shape together determine
its output shape. In a concatenation expression L R, the output shape of L is
the input shape of R.
*/
package shapes

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type Shape struct {
	Type  types.Type
	Stack *Stack
}

type Stack struct {
	Head NFF
	Tail *Stack
}

func (stack *Stack) Push(n NFF) *Stack {
	return &Stack{n, stack}
}

func (stack *Stack) Pop() *Stack {
	return stack.Tail
}

type NFF struct {
	InputType types.Type // TODO type parameters
	Name      string     // TODO namespaces
	ArgTypes  []types.Type
	Funcer    func([]Function) Function // TODO first-class functions
}

type Function interface {
	OutputShape(inputShape Shape) Shape
	OutputState(inputState states.State) states.State
}

func (inputShape Shape) ResolveNFF(Pos lexer.Position, name string, argFunctions []Function) (Function, error) {
	argShapes := make([]Shape, len(argFunctions))
	for i, f := range argFunctions {
		argShapes[i] = f.OutputShape(Shape{&types.AnyType{}, inputShape.Stack})
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
