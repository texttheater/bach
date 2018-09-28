package ast

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/contexts"
	"github.com/texttheater/bach/functions"
)

///////////////////////////////////////////////////////////////////////////////

type Expression interface {
	Function(context contexts.Context) functions.Function
}

///////////////////////////////////////////////////////////////////////////////

type IdentityExpression struct {
}

func (x IdentityExpression) Function(context contexts.Context) functions.Function {
	return functions.IdentityFunction{context.InputType}
}

///////////////////////////////////////////////////////////////////////////////

type CompositionExpression struct {
	Left Expression
	Right Expression
}

func (x CompositionExpression) Function(leftContext contexts.Context) functions.Function {
	leftFunction := x.Left.Function(leftContext)
	middleContext := contexts.Context{leftFunction.Type()}
	rightFunction := x.Right.Function(middleContext)
	return functions.CompositionFunction{leftFunction, rightFunction}
}

///////////////////////////////////////////////////////////////////////////////

type NumberExpression struct {
	Pos lexer.Position
	Value float64
}

func (x NumberExpression) Function(context contexts.Context) functions.Function {
	return functions.NumberFunction{x.Value}
}
