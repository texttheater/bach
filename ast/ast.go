package ast

import (
	//"fmt"
	//"os"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/nffs"
	"github.com/texttheater/bach/types"
)

///////////////////////////////////////////////////////////////////////////////

type Expression interface {
	Function(inputType types.Type) (functions.Function, error)
}

///////////////////////////////////////////////////////////////////////////////

type IdentityExpression struct {
}

func (x *IdentityExpression) Function(inputType types.Type) (functions.Function, error) {
	return &functions.IdentityFunction{inputType}, nil
}

///////////////////////////////////////////////////////////////////////////////

type CompositionExpression struct {
	Left  Expression
	Right Expression
}

func (x *CompositionExpression) Function(inputType types.Type) (functions.Function, error) {
	leftFunction, err := x.Left.Function(inputType)
	if err != nil {
		return nil, err
	}
	rightFunction, err := x.Right.Function(leftFunction.Type())
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

func (x *NumberExpression) Function(inputType types.Type) (functions.Function, error) {
	return &functions.NumberFunction{x.Value}, nil
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

func (x *NFFCallExpression) Function(inputType types.Type) (functions.Function, error) {
	argFunctions := make([]functions.Function, len(x.Args))
	for i, arg := range x.Args {
		f, err := arg.Function(&types.AnyType{})
		if err != nil {
			return nil, err
		}
		argFunctions[i] = f
	}
	f, err := nffs.Function(x.Pos, inputType, x.Name, argFunctions)
	if err != nil {
		return nil, err
	}
	return f, nil
}
