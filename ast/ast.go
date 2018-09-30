package ast

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/contexts"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/nffs"
)

///////////////////////////////////////////////////////////////////////////////

type Expression interface {
	Function(context contexts.Context) (functions.Function, error)
}

///////////////////////////////////////////////////////////////////////////////

type IdentityExpression struct {
}

func (x IdentityExpression) Function(context contexts.Context) (functions.Function, error) {
	return functions.IdentityFunction{context.InputType}, nil
}

///////////////////////////////////////////////////////////////////////////////

type CompositionExpression struct {
	Left Expression
	Right Expression
}

func (x CompositionExpression) Function(leftContext contexts.Context) (functions.Function, error) {
	leftFunction, err := x.Left.Function(leftContext)
	if err != nil {
		return nil, err
	}
	middleContext := contexts.Context{leftFunction.Type()}
	rightFunction, err := x.Right.Function(middleContext)
	if err != nil {
		return nil, err
	}
	return functions.CompositionFunction{leftFunction, rightFunction}, nil
}

///////////////////////////////////////////////////////////////////////////////

type NumberExpression struct {
	Pos lexer.Position
	Value float64
}

func (x NumberExpression) Function(context contexts.Context) (functions.Function, error) {
	return functions.NumberFunction{x.Value}, nil
}

///////////////////////////////////////////////////////////////////////////////

// NFF = named function family (close to what is called a function in most
// programming languages). TODO rename to something prettier.

// TODO namespaces

type NFFCallExpression struct {
	Pos lexer.Position
	Name string
	Args []Expression
}

func (x NFFCallExpression) Function(context contexts.Context) (functions.Function, error) {
	argFunctions := make([]functions.Function, len(x.Args))
	for i, arg := range x.Args {
		f, err := arg.Function(context) // TODO allow other contexts
		if err != nil {
			return nil, err
		}
		argFunctions[i] = f
	}
	f, err := nffs.Function(x.Pos, context.InputType, x.Name, argFunctions)
	if err != nil {
		return nil, err
	}
	return f, nil
}
