/*
Package ast implements Bach's abstract syntax trees.

An alternative name for this package would be: expressions. Because that's what
an AST is, an expression consisting of subexpressions.
*/
package ast

import (
	"fmt"
	//"os"
	"strings"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

///////////////////////////////////////////////////////////////////////////////

type Expression interface {
	Function(inputShape functions.Shape, params []*functions.Parameter) (*functions.Function, error)
}

///////////////////////////////////////////////////////////////////////////////

type CompositionExpression struct {
	Pos   lexer.Position
	Left  Expression
	Right Expression
}

func (x *CompositionExpression) Function(inputShape functions.Shape, params []*functions.Parameter) (*functions.Function, error) {
	if len(params) > 0 {
		return nil, errors.E("type", x.Pos, fmt.Sprintf("%s parameters required here, but composition expressions have no parameters", len(params)))
	}
	leftFunction, err := x.Left.Function(inputShape, nil)
	if err != nil {
		return nil, err
	}
	middleShape := leftFunction.UpdateShape(inputShape)
	rightFunction, err := x.Right.Function(middleShape, nil)
	if err != nil {
		return nil, err
	}
	return &functions.Function{
		inputShape.Type,
		"",
		nil,
		nil,
		func(inputShape functions.Shape) functions.Shape {
			return rightFunction.UpdateShape(middleShape)
		},
		func(inputState functions.State, args []*functions.Function) functions.State {
			middleState := leftFunction.Apply(inputState, nil)
			return rightFunction.Apply(middleState, nil)
		},
	}, nil
}

///////////////////////////////////////////////////////////////////////////////

type NumberExpression struct {
	Pos   lexer.Position
	Value float64
}

func (x *NumberExpression) Function(inputShape functions.Shape, params []*functions.Parameter) (*functions.Function, error) {
	if len(params) > 0 {
		return nil, errors.E("type", x.Pos, fmt.Sprintf("%s parameters required here, but number expressions have no parameters", len(params)))
	}
	return functions.SimpleFunction(
		&types.AnyType{},
		"",
		nil,
		&types.NumberType{},
		func(values.Value, []values.Value) values.Value {
			return &values.NumberValue{x.Value}
		},
	), nil
}

///////////////////////////////////////////////////////////////////////////////

type StringExpression struct {
	Pos   lexer.Position
	Value string
}

func (x *StringExpression) Function(inputShape functions.Shape, params []*functions.Parameter) (*functions.Function, error) {
	if len(params) > 0 {
		return nil, errors.E("type", x.Pos, fmt.Sprintf("%s parameters required here, but string expressions have no parameters", len(params)))
	}
	return functions.SimpleFunction(
		&types.AnyType{},
		"",
		nil,
		&types.StringType{},
		func(values.Value, []values.Value) values.Value {
			return &values.StringValue{x.Value}
		},
	), nil
}

///////////////////////////////////////////////////////////////////////////////

type CallExpression struct {
	Pos       lexer.Position
	Name      string // TODO namespaces
	Arguments []Expression
}

func (x *CallExpression) Function(inputShape functions.Shape, params []*functions.Parameter) (*functions.Function, error) {
	stack := inputShape.Stack
	for stack != nil {
		function := stack.Head
		if function.Name != x.Name {
			stack = stack.Tail
			continue
		}
		if len(function.OpenParams) != len(x.Arguments)+len(params) {
			stack = stack.Tail
			continue
		}
		if !function.InputType.Subsumes(inputShape.Type) {
			stack = stack.Tail
			continue
		}
		for i, arg := range x.Arguments {
			argInputShape := functions.Shape{
				function.OpenParams[0].InputType,
				inputShape.Stack,
			}
			argFunction, err := arg.Function(argInputShape,
				function.OpenParams[0].Parameters)
			if err != nil {
				return nil, err
			}
			argOutputType := argFunction.UpdateShape(
				argFunction.UpdateShape(argInputShape),
			).Type
			if !function.OpenParams[0].OutputType.Subsumes(argOutputType) {
				return nil, errors.E("type", x.Pos,
					"argument #%s needs output type %s, got %s", i,
					function.OpenParams[0],
					argOutputType)
			}
			function = function.SetArg(argFunction)
		}
		for i, param := range params {
			if !param.Subsumes(function.OpenParams[i]) {
				return nil, errors.E("type", x.Pos,
					"parameter #%s must be %s, got %s", i,
					param, function.OpenParams[i])
			}
		}
		return function, nil
	}
	return nil, errors.E("type", x.Pos, "no such function")
}

func formatArgTypes(argShapes []functions.Shape) string {
	formatted := make([]string, len(argShapes))
	for i, s := range argShapes {
		formatted[i] = fmt.Sprintf("%v", s.Type)
	}
	return strings.Join(formatted, ", ")
}

///////////////////////////////////////////////////////////////////////////////

type AssignmentExpression struct {
	Pos  lexer.Position
	Name string
}

func (x *AssignmentExpression) Function(inputShape functions.Shape, params []*functions.Parameter) (*functions.Function, error) {
	if len(params) > 0 {
		return nil, errors.E("type", x.Pos, fmt.Sprintf("%s parameters required here, but assignment expressions have no parameters", len(params)))
	}
	return &functions.Function{
		inputShape.Type,
		"",
		nil,
		nil,
		func(inputShape functions.Shape) functions.Shape {
			return functions.Shape{
				inputShape.Type,
				inputShape.Stack.Push(&functions.Function{
					&types.AnyType{},
					x.Name,
					nil,
					nil,
					func(varInputShape functions.Shape) functions.Shape {
						return functions.Shape{
							inputShape.Type,
							varInputShape.Stack,
						}
					},
					func(inputState functions.State, argFunctions []*functions.Function) functions.State {
						stack := inputState.Stack
						for stack != nil {
							if stack.Head.Name == x.Name {
								return functions.State{
									stack.Head.Value,
									inputState.Stack,
								}
							}
							stack = stack.Tail
						}
						panic("unknown variable")
					},
				}),
			}
		},
		func(inputState functions.State, argFunctions []*functions.Function) functions.State {
			return functions.State{
				inputState.Value,
				inputState.Stack.Push(functions.NamedValue{
					x.Name,
					inputState.Value,
				}),
			}
		},
	}, nil
}

///////////////////////////////////////////////////////////////////////////////
