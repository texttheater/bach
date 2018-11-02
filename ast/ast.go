/*
Package ast implements Bach's abstract syntax trees.

An alternative name for this package would be: expressions. Because that's what
an AST is, an expression consisting of subexpressions.
*/
package ast

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

///////////////////////////////////////////////////////////////////////////////

var nullContext = functions.Context{}

///////////////////////////////////////////////////////////////////////////////

type Expression interface {
	Typecheck(inputContext functions.Context, params []*functions.Param) (outputContext functions.Context, action *functions.Action, err error)
}

///////////////////////////////////////////////////////////////////////////////

type ConstantExpression struct {
	Pos   lexer.Position
	Type  types.Type
	Value values.Value
}

func (x *ConstantExpression) Typecheck(inputContext functions.Context, params []*functions.Param) (functions.Context, *functions.Action, error) {
	if len(params) > 0 {
		return nullContext, nil, errors.E("type", x.Pos, "number expression does not take parameters")
	}
	outputContext := functions.Context{x.Type, inputContext.FunctionStack}
	action := &functions.Action{
		Execute: func(inputValue values.Value, args []*functions.Action) values.Value {
			return x.Value
		},
	}
	return outputContext, action, nil
}

///////////////////////////////////////////////////////////////////////////////

type CompositionExpression struct {
	Pos   lexer.Position
	Left  Expression
	Right Expression
}

func (x *CompositionExpression) Typecheck(inputContext functions.Context, params []*functions.Param) (functions.Context, *functions.Action, error) {
	if len(params) > 0 {
		return nullContext, nil, errors.E("type", x.Pos, "composition expression does not take parameters")
	}
	middleContext, lAction, err := x.Left.Typecheck(inputContext, nil)
	if err != nil {
		return nullContext, nil, err
	}
	outputContext, rAction, err := x.Right.Typecheck(middleContext, nil)
	if err != nil {
		return nullContext, nil, err
	}
	action := &functions.Action{
		Execute: func(inputValue values.Value, args []*functions.Action) values.Value {
			middleValue := lAction.Execute(inputValue, nil)
			outputValue := rAction.Execute(middleValue, nil)
			return outputValue
		},
	}
	return outputContext, action, nil
}

///////////////////////////////////////////////////////////////////////////////

type CallExpression struct {
	Pos  lexer.Position
	Name string
	Args []Expression
}

func (x *CallExpression) Typecheck(inputContext functions.Context, params []*functions.Param) (functions.Context, *functions.Action, error) {
	// Go down the function stack and find the function invoked by this
	// call
	stack := inputContext.FunctionStack
Entries:
	for {
		// Reached bottom of stack without finding a matching function
		if stack == nil {
			return nullContext, nil, errors.E("type", x.Pos, "no such function")
		}
		// Try the function on top of the stack
		function := stack.Head
		// Check function name
		if function.Name != x.Name {
			stack = stack.Tail
			continue
		}
		// Check number of params
		if len(function.Params) != len(x.Args)+len(params) {
			stack = stack.Tail
			continue
		}
		// Check input type
		if !function.InputType.Subsumes(inputContext.Type) {
			stack = stack.Tail
			continue
		}
		// Prepare action:
		action := function.Action
		// Check function params filled by this call
		for i := 0; i < len(x.Args); i++ {
			param := function.Params[i]
			argInputContext := functions.Context{param.InputType, inputContext.FunctionStack}
			argOutputContext, argAction, err := x.Args[i].Typecheck(argInputContext, param.Params)
			if err != nil {
				stack = stack.Tail
				continue Entries
			}
			if !param.OutputType.Subsumes(argOutputContext.Type) {
				stack = stack.Tail
				continue Entries
			}
			action = action.SetArg(argAction)
		}
		// Check function params *not* filled by this call (thus left
		// to function to call)
		for i := 0; i < len(params); i++ {
			if !params[i].Subsumes(function.Params[len(x.Args)+i]) {
				stack = stack.Tail
				continue Entries
			}
		}
		// Return result
		return functions.Context{function.OutputType, inputContext.FunctionStack}, action, nil
	}
}

///////////////////////////////////////////////////////////////////////////////

type AssignmentExpression struct {
	Pos  lexer.Position
	Name string
}

func (x *AssignmentExpression) Typecheck(inputContext functions.Context, params []*functions.Param) (functions.Context, *functions.Action, error) {
	var value values.Value
	outputContext := functions.Context{inputContext.Type, inputContext.FunctionStack.Push(functions.Function{
		InputType:  &types.AnyType{},
		Name:       x.Name,
		Params:     nil,
		OutputType: inputContext.Type,
		Action: &functions.Action{
			Execute: func(inputValue values.Value, args []*functions.Action) values.Value {
				return value
			},
		},
	})}
	action := &functions.Action{
		Execute: func(inputValue values.Value, args []*functions.Action) values.Value {
			value = inputValue
			return inputValue
		},
	}
	return outputContext, action, nil
}

///////////////////////////////////////////////////////////////////////////////
