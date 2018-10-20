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
	Funcer(inputShape functions.Shape) (functions.Funcer, error)
}

///////////////////////////////////////////////////////////////////////////////

type IdentityExpression struct {
}

func (x *IdentityExpression) Funcer(inputShape functions.Shape) (functions.Funcer, error) {
	return functions.NoArgFuncer(&functions.IdentityFunction{inputShape.Type}), nil
}

///////////////////////////////////////////////////////////////////////////////

type CompositionExpression struct {
	Left  Expression
	Right Expression
}

func (x *CompositionExpression) Funcer(inputShape functions.Shape) (functions.Funcer, error) {
	leftFuncer, err := x.Left.Funcer(inputShape)
	if err != nil {
		return nil, err
	}
	leftFunction := leftFuncer()
	rightFuncer, err := x.Right.Funcer(leftFuncer().OutputShape(inputShape.Stack))
	if err != nil {
		return nil, err
	}
	return functions.NoArgFuncer(&functions.CompositionFunction{leftFunction, rightFuncer()}), nil
}

///////////////////////////////////////////////////////////////////////////////

type NumberExpression struct {
	Pos   lexer.Position
	Value float64
}

func (x *NumberExpression) Funcer(inputShape functions.Shape) (functions.Funcer, error) {
	return functions.NoArgFuncer(&functions.LiteralFunction{&types.NumberType{}, &values.NumberValue{x.Value}}), nil
}

///////////////////////////////////////////////////////////////////////////////

type StringExpression struct {
	Pos   lexer.Position
	Value string
}

func (x *StringExpression) Funcer(inputShape functions.Shape) (functions.Funcer, error) {
	return functions.NoArgFuncer(&functions.LiteralFunction{&types.StringType{}, &values.StringValue{x.Value}}), nil
}

///////////////////////////////////////////////////////////////////////////////

// NFF = named function family (close to what is called a function in most
// programming languages). TODO rename to something prettier.

// TODO namespaces

type NFFCallExpression struct {
	Pos       lexer.Position
	Name      string
	Arguments []Expression
}

func (x *NFFCallExpression) Funcer(inputShape functions.Shape) (functions.Funcer, error) {
	stack := inputShape.Stack
Entries:
	for stack != nil {
		if stack.Head.Name != x.Name {
			stack = stack.Tail
			continue
		}
		if len(stack.Head.Parameters) != len(x.Arguments) {
			stack = stack.Tail
			continue
		}
		if !stack.Head.InputType.Subsumes(inputShape.Type) {
			stack = stack.Tail
			continue
		}
		argFuncers := make([]functions.Funcer, 0, len(x.Arguments))
		for i, par := range stack.Head.Parameters {
			argFuncer, err := x.Arguments[i].Funcer(functions.Shape{par.InputType, inputShape.Stack})
			if err != nil {
				continue Entries
			}
			argFunction := argFuncer() // TODO pass arguments
			argOutputType := argFunction.OutputShape(inputShape.Stack).Type
			if !par.OutputType.Subsumes(argOutputType) {
				continue Entries
			}
			argFuncers = append(argFuncers, argFuncer)
		}
		return func(preFuncers ...functions.Funcer) functions.Function {
			for _, argFuncer := range argFuncers {
				preFuncers = append(preFuncers, argFuncer)
			}
			return stack.Head.Funcer(preFuncers...)
		}, nil
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

func (x *AssignmentExpression) Funcer(inputShape functions.Shape) (functions.Funcer, error) {
	return functions.NoArgFuncer(&functions.AssignmentFunction{inputShape.Type, x.Name}), nil
}

///////////////////////////////////////////////////////////////////////////////
