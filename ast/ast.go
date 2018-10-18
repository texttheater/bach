/*
Package ast implements Bach's abstract syntax trees.

An alternative name for this package would be: expressions. Because that's what
an AST is, an expression consisting of subexpressions.
*/
package ast

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

///////////////////////////////////////////////////////////////////////////////

type Expression interface {
	Function(inputShape functions.Shape) (functions.Function, error)
}

///////////////////////////////////////////////////////////////////////////////

type IdentityExpression struct {
}

func (x *IdentityExpression) Function(inputShape functions.Shape) (functions.Function, error) {
	return &functions.IdentityFunction{inputShape.Type}, nil
}

///////////////////////////////////////////////////////////////////////////////

type CompositionExpression struct {
	Left  Expression
	Right Expression
}

func (x *CompositionExpression) Function(inputShape functions.Shape) (functions.Function, error) {
	leftFunction, err := x.Left.Function(inputShape)
	if err != nil {
		return nil, err
	}
	rightFunction, err := x.Right.Function(leftFunction.OutputShape(inputShape.Stack))
	if err != nil {
		return nil, err
	}
	return &functions.CompositionFunction{leftFunction, rightFunction}, nil
}

///////////////////////////////////////////////////////////////////////////////

type NumberExpression struct {
	Pos   lexer.Position
	Value float64
}

func (x *NumberExpression) Function(inputShape functions.Shape) (functions.Function, error) {
	return &functions.LiteralFunction{&types.NumberType{}, &values.NumberValue{x.Value}}, nil
}

///////////////////////////////////////////////////////////////////////////////

type StringExpression struct {
	Pos   lexer.Position
	Value string
}

func (x *StringExpression) Function(inputShape functions.Shape) (functions.Function, error) {
	return &functions.LiteralFunction{&types.StringType{}, &values.StringValue{x.Value}}, nil
}

///////////////////////////////////////////////////////////////////////////////

// NFF = named function family (close to what is called a function in most
// programming languages). TODO rename to something prettier.

// TODO namespaces

type NFFCallExpression struct {
	Pos  lexer.Position
	Name string
	Args []Expression
}

func (x *NFFCallExpression) Function(inputShape functions.Shape) (functions.Function, error) {
	// Step 1: resolve arguments to functions
	argFunctions := make([]functions.Function, len(x.Args))
	for i, arg := range x.Args {
		f, err := arg.Function(inputShape)
		if err != nil {
			return nil, err
		}
		argFunctions[i] = f
	}
	// Step 2: search the stack for matching NFFs
	stack := inputShape.Stack
	for stack != nil {
		function, ok := stack.Head.Function(inputShape, x.Name, argFunctions)
		if ok {
			return function, nil
		}
		stack = stack.Tail
	}
	// Fail:
	argShapes := make([]functions.Shape, len(argFunctions))
	for i, f := range argFunctions {
		argShapes[i] = f.OutputShape(inputShape.Stack)
	}
	return nil, errors.E("type", x.Pos, "no function found: for %v %v(%s)", inputShape.Type, x.Name, formatArgTypes(argShapes))
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

func (x *AssignmentExpression) Function(inputShape functions.Shape) (functions.Function, error) {
	return &functions.AssignmentFunction{inputShape.Type, x.Name}, nil
}

///////////////////////////////////////////////////////////////////////////////
