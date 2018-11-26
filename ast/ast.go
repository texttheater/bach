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

var booleanType = types.BooleanType{}

///////////////////////////////////////////////////////////////////////////////

type Expression interface {
	Typecheck(inputContext functions.Context, params []*functions.Param) (outputContext functions.Context, action functions.Action, err error)
}

///////////////////////////////////////////////////////////////////////////////

type ConstantExpression struct {
	Pos   lexer.Position
	Type  types.Type
	Value values.Value
}

func (x *ConstantExpression) Typecheck(inputContext functions.Context, params []*functions.Param) (functions.Context, functions.Action, error) {
	if len(params) > 0 {
		return nullContext, nil, errors.E("type", x.Pos, "number expression does not take parameters")
	}
	outputContext := functions.Context{x.Type, inputContext.FunctionStack}
	action := func(inputValue values.Value, args []functions.Action) values.Value {
		return x.Value
	}
	return outputContext, action, nil
}

///////////////////////////////////////////////////////////////////////////////

type CompositionExpression struct {
	Pos   lexer.Position
	Left  Expression
	Right Expression
}

func (x *CompositionExpression) Typecheck(inputContext functions.Context, params []*functions.Param) (functions.Context, functions.Action, error) {
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
	action := func(inputValue values.Value, args []functions.Action) values.Value {
		middleValue := lAction(inputValue, nil)
		outputValue := rAction(middleValue, nil)
		return outputValue
	}
	return outputContext, action, nil
}

///////////////////////////////////////////////////////////////////////////////

type CallExpression struct {
	Pos  lexer.Position
	Name string
	Args []Expression
}

func (x *CallExpression) Typecheck(inputContext functions.Context, params []*functions.Param) (functions.Context, functions.Action, error) {
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

func (x *AssignmentExpression) Typecheck(inputContext functions.Context, params []*functions.Param) (functions.Context, functions.Action, error) {
	if len(params) > 0 {
		return nullContext, nil, errors.E("type", x.Pos, "assignment expression does not take parameters")
	}
	var value values.Value
	outputContext := functions.Context{inputContext.Type, inputContext.FunctionStack.Push(functions.Function{
		InputType:  &types.AnyType{},
		Name:       x.Name,
		Params:     nil,
		OutputType: inputContext.Type,
		Action: func(inputValue values.Value, args []functions.Action) values.Value {
			return value
		},
	})}
	action := func(inputValue values.Value, args []functions.Action) values.Value {
		value = inputValue
		return inputValue
	}
	return outputContext, action, nil
}

///////////////////////////////////////////////////////////////////////////////

type DefinitionExpression struct {
	Pos        lexer.Position
	InputType  types.Type
	Name       string
	Params     []*functions.Param
	OutputType types.Type
	Body       Expression
}

func (x *DefinitionExpression) Typecheck(inputContext functions.Context, params []*functions.Param) (functions.Context, functions.Action, error) {
	// make sure we got no parameters
	if len(params) > 0 {
		return nullContext, nil, errors.E("type", x.Pos, "definition expression does not take parameters")
	}
	// variable for body action (will be set later)
	var bodyAction functions.Action = nil
	// add the function defined here to the function stack
	functionStack := inputContext.FunctionStack.Push(functions.Function{
		InputType: x.InputType,
		Name: x.Name,
		Params: x.Params,
		OutputType: x.OutputType,
		Action: func(inputValue values.Value, args []functions.Action) values.Value {
			// Push, call, pop
			for i, param := range x.Params {
				param.ActionStack = param.ActionStack.Push(args[i])
			}
			result := bodyAction(inputValue, nil)
			for _, param := range x.Params {
				param.ActionStack = param.ActionStack.Tail
			}
			return result
		},
	})
	// add parameter functions for use in the body
	bodyFunctionStack := functionStack
	for _, param := range x.Params {
		bodyFunctionStack = bodyFunctionStack.Push(functions.Function{
			InputType:  param.InputType,
			Name:       param.Name,
			Params:     param.Params,
			OutputType: param.OutputType,
			Action: func(inputValue values.Value, args []functions.Action) values.Value {
				return param.ActionStack.Head(inputValue, args)
			},
		})
	}
	// define body input context
	bodyInputContext := functions.Context{
		Type:          x.InputType,
		FunctionStack: bodyFunctionStack,
	}
	// typecheck body (crucially, setting body action)
	bodyOutputContext, bodyAction, err := x.Body.Typecheck(bodyInputContext, nil)
	if err != nil {
		return nullContext, nil, err
	}
	// check body output type
	if !x.OutputType.Subsumes(bodyOutputContext.Type) {
		return nullContext, nil, errors.E("type", x.Pos, "expected function body output type %s, got %s", x.OutputType, bodyOutputContext.Type)
	}
	// define output context
	outputContext := functions.Context{
		Type: inputContext.Type,
		FunctionStack: functionStack,
	}
	// define action (simple identity)
	action := func(inputValue values.Value, args []functions.Action) values.Value {
		return inputValue
	}
	// return
	return outputContext, action, nil
}

///////////////////////////////////////////////////////////////////////////////

type ConditionalExpression struct {
	Pos             lexer.Position
	Condition       Expression
	Consequent      Expression
	ElifConditions  []Expression
	ElifConsequents []Expression
	Alternative     Expression
}

func (x *ConditionalExpression) Typecheck(inputContext functions.Context, params []*functions.Param) (functions.Context, functions.Action, error) {
	conditionOutputContext, conditionAction, err := x.Condition.Typecheck(inputContext, nil)
	if err != nil {
		return nullContext, nil, err
	}
	if !booleanType.Subsumes(conditionOutputContext.Type) {
		return nullContext, nil, errors.E("type", x.Pos, "condition must be boolean")
	}
	// context is the shared input context for all conditions and consequents.
	// Each condition may add to the FunctionStack. Type always stays the same.
	context := functions.Context{
		Type:          inputContext.Type,
		FunctionStack: conditionOutputContext.FunctionStack,
	}
	consequentOutputContext, consequentAction, err := x.Consequent.Typecheck(context, nil)
	if err != nil {
		return nullContext, nil, err
	}
	outputType := consequentOutputContext.Type
	elifConditionActions := make([]functions.Action, 0, len(x.ElifConditions))
	elifConsequentActions := make([]functions.Action, 0, len(x.ElifConsequents))
	for i := range x.ElifConditions {
		conditionOutputContext, elifConditionAction, err := x.ElifConditions[i].Typecheck(context, nil)
		if err != nil {
			return nullContext, nil, err
		}
		if !booleanType.Subsumes(conditionOutputContext.Type) {
			return nullContext, nil, errors.E("type", x.Pos, "condition must be boolean")
		}
		context.FunctionStack = conditionOutputContext.FunctionStack
		elifConditionActions = append(elifConditionActions, elifConditionAction)
		consequentOutputContext, elifConsequentAction, err := x.ElifConsequents[i].Typecheck(context, nil)
		if err != nil {
			return nullContext, nil, err
		}
		elifConsequentActions = append(elifConsequentActions, elifConsequentAction)
		outputType = types.Disjoin(outputType, consequentOutputContext.Type)
	}
	alternativeOutputContext, alternativeAction, err := x.Alternative.Typecheck(context, nil)
	if err != nil {
		return nullContext, nil, err
	}
	outputType = types.Disjoin(outputType, alternativeOutputContext.Type)
	action := func(inputValue values.Value, args []functions.Action) values.Value {
		conditionValue := conditionAction(inputValue, nil)
		boolConditionValue, _ := conditionValue.(*values.BooleanValue)
		if boolConditionValue.Value {
			return consequentAction(inputValue, nil)
		}
		for i := range elifConditionActions {
			conditionValue = elifConditionActions[i](inputValue, nil)
			boolConditionValue, _ = conditionValue.(*values.BooleanValue)
			if boolConditionValue.Value {
				return elifConsequentActions[i](inputValue, nil)
			}
		}
		return alternativeAction(inputValue, nil)
	}
	outputContext := functions.Context{
		Type:          outputType,
		FunctionStack: inputContext.FunctionStack,
	}
	return outputContext, action, nil
}

///////////////////////////////////////////////////////////////////////////////
