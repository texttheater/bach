/*
Package ast implements Bach's abstract syntax trees.

An alternative name for this package would be: expressions. Because that's what
an AST is, an expression consisting of subexpressions.
*/
package ast

import (
	//"fmt"
	//"os"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

///////////////////////////////////////////////////////////////////////////////

type Expression interface {
	Function(inputShape shapes.Shape) (shapes.Function, error)
}

///////////////////////////////////////////////////////////////////////////////

type IdentityExpression struct {
}

func (x *IdentityExpression) Function(inputShape shapes.Shape) (shapes.Function, error) {
	return &functions.IdentityFunction{}, nil
}

///////////////////////////////////////////////////////////////////////////////

type CompositionExpression struct {
	Left  Expression
	Right Expression
}

func (x *CompositionExpression) Function(inputShape shapes.Shape) (shapes.Function, error) {
	leftFunction, err := x.Left.Function(inputShape)
	if err != nil {
		return nil, err
	}
	rightFunction, err := x.Right.Function(leftFunction.OutputShape(inputShape))
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

func (x *NumberExpression) Function(inputShape shapes.Shape) (shapes.Function, error) {
	return &functions.LiteralFunction{&types.NumberType{}, &values.NumberValue{x.Value}}, nil
}

///////////////////////////////////////////////////////////////////////////////

type StringExpression struct {
	Pos   lexer.Position
	Value string
}

func (x *StringExpression) Function(inputShape shapes.Shape) (shapes.Function, error) {
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

func (x *NFFCallExpression) Function(inputShape shapes.Shape) (shapes.Function, error) {
	argFunctions := make([]shapes.Function, len(x.Args))
	for i, arg := range x.Args {
		f, err := arg.Function(inputShape)
		if err != nil {
			return nil, err
		}
		argFunctions[i] = f
	}
	f, err := inputShape.ResolveNFF(x.Pos, x.Name, argFunctions)
	if err != nil {
		return nil, err
	}
	return f, nil
}

///////////////////////////////////////////////////////////////////////////////

type AssignmentExpression struct {
	Pos  lexer.Position
	Name string
}

func (x *AssignmentExpression) Function(inputShape shapes.Shape) (shapes.Function, error) {
	return &functions.AssignmentFunction{x.Name}, nil
}

///////////////////////////////////////////////////////////////////////////////
