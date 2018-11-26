/*
Package ast implements Bach's abstract syntax trees.

An alternative name for this package would be: expressions. Because that's what
an AST is, an expression consisting of subexpressions.
*/
package ast

import (
	"fmt"
	//"os"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

///////////////////////////////////////////////////////////////////////////////

var nullShape = shapes.Shape{}

var booleanType = types.BooleanType{}

///////////////////////////////////////////////////////////////////////////////

type Expression interface {
	Typecheck(inputShape shapes.Shape, params []*functions.Param) (outputShape shapes.Shape, action functions.Action, err error)
}

///////////////////////////////////////////////////////////////////////////////

type ConstantExpression struct {
	Pos   lexer.Position
	Type  types.Type
	Value values.Value
}

func (x *ConstantExpression) Typecheck(inputShape shapes.Shape, params []*functions.Param) (shapes.Shape, functions.Action, error) {
	if len(params) > 0 {
		return nullShape, nil, errors.E("type", x.Pos, "number expression does not take parameters")
	}
	outputShape := shapes.Shape{x.Type, inputShape.FunctionStack}
	action := func(inputState states.State, args []functions.Action) states.State {
		return states.State{
			Value: x.Value,
			Stack: inputState.Stack,
		}
	}
	return outputShape, action, nil
}

///////////////////////////////////////////////////////////////////////////////

type CompositionExpression struct {
	Pos   lexer.Position
	Left  Expression
	Right Expression
}

func (x *CompositionExpression) Typecheck(inputShape shapes.Shape, params []*functions.Param) (shapes.Shape, functions.Action, error) {
	if len(params) > 0 {
		return nullShape, nil, errors.E("type", x.Pos, "composition expression does not take parameters")
	}
	middleShape, lAction, err := x.Left.Typecheck(inputShape, nil)
	if err != nil {
		return nullShape, nil, err
	}
	outputShape, rAction, err := x.Right.Typecheck(middleShape, nil)
	if err != nil {
		return nullShape, nil, err
	}
	action := func(inputState states.State, args []functions.Action) states.State {
		middleState := lAction(inputState, nil)
		outputState := rAction(middleState, nil)
		return outputState
	}
	return outputShape, action, nil
}

///////////////////////////////////////////////////////////////////////////////

type CallExpression struct {
	Pos  lexer.Position
	Name string
	Args []Expression
}

func (x *CallExpression) Typecheck(inputShape shapes.Shape, params []*functions.Param) (shapes.Shape, functions.Action, error) {
	// Go down the function stack and find the function invoked by this
	// call
	stack := inputShape.FunctionStack
Entries:
	for {
		// Reached bottom of stack without finding a matching function
		if stack == nil {
			return nullShape, nil, errors.E("type", x.Pos, "no such function")
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
		if !function.InputType.Subsumes(inputShape.Type) {
			stack = stack.Tail
			continue
		}
		// Prepare action:
		action := function.Action
		// Check function params filled by this call
		for i := 0; i < len(x.Args); i++ {
			param := function.Params[i]
			argInputShape := shapes.Shape{param.InputType, inputShape.FunctionStack}
			argOutputShape, argAction, err := x.Args[i].Typecheck(argInputShape, param.Params)
			if err != nil {
				stack = stack.Tail
				continue Entries
			}
			if !param.OutputType.Subsumes(argOutputShape.Type) {
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
		return shapes.Shape{function.OutputType, inputShape.FunctionStack}, action, nil
	}
}

///////////////////////////////////////////////////////////////////////////////

type AssignmentExpression struct {
	Pos  lexer.Position
	Name string
}

func (x *AssignmentExpression) Typecheck(inputShape shapes.Shape, params []*functions.Param) (shapes.Shape, functions.Action, error) {
	if len(params) > 0 {
		return nullShape, nil, errors.E("type", x.Pos, "assignment expression does not take parameters")
	}
	outputShape := shapes.Shape{inputShape.Type, inputShape.FunctionStack.Push(functions.Function{
		InputType:  &types.AnyType{},
		Name:       x.Name,
		Params:     nil,
		OutputType: inputShape.Type,
		Action: func(inputState states.State, args []functions.Action) states.State {
			stack := inputState.Stack
			for stack != nil {
				if stack.Head.Name == x.Name {
					return states.State{
						Value: stack.Head.Value,
						Stack: inputState.Stack,
					}
				}
				stack = stack.Tail
			}
			panic(fmt.Sprintf("variable %s not found", x.Name))
		},
	})}
	action := func(inputState states.State, args []functions.Action) states.State {
		return states.State{
			Value: inputState.Value,
			Stack: inputState.Stack.Push(states.Variable{
				Name:  x.Name,
				Value: inputState.Value,
			}),
		}
	}
	return outputShape, action, nil
}

type valueStack struct {
	Head values.Value
	Tail *valueStack
}

func (s *valueStack) Push(element values.Value) *valueStack {
	return &valueStack{element, s}
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

func (x *DefinitionExpression) Typecheck(inputShape shapes.Shape, params []*functions.Param) (shapes.Shape, functions.Action, error) {
	// make sure we got no parameters
	if len(params) > 0 {
		return nullShape, nil, errors.E("type", x.Pos, "definition expression does not take parameters")
	}
	// variables for body input stack, action (will be set later)
	var bodyInputStack *states.Stack = nil
	var bodyAction functions.Action = nil
	// add the function defined here to the function stack
	functionStack := inputShape.FunctionStack.Push(functions.Function{
		InputType:  x.InputType,
		Name:       x.Name,
		Params:     x.Params,
		OutputType: x.OutputType,
		Action: func(inputState states.State, args []functions.Action) states.State {
			// Push, call, pop
			for i, param := range x.Params {
				param.ActionStack = param.ActionStack.Push(args[i])
			}
			bodyInputState := states.State{
				Value: inputState.Value,
				Stack: bodyInputStack,
			}
			bodyOutputState := bodyAction(bodyInputState, nil)
			for _, param := range x.Params {
				param.ActionStack = param.ActionStack.Tail
			}
			return states.State{
				Value: bodyOutputState.Value,
				Stack: inputState.Stack,
			}
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
			Action: func(inputState states.State, args []functions.Action) states.State {
				return param.ActionStack.Head(inputState, args)
			},
		})
	}
	// define body input context
	bodyInputShape := shapes.Shape{
		Type:          x.InputType,
		FunctionStack: bodyFunctionStack,
	}
	// typecheck body (crucially, setting body action)
	bodyOutputShape, bodyAction, err := x.Body.Typecheck(bodyInputShape, nil)
	if err != nil {
		return nullShape, nil, err
	}
	// check body output type
	if !x.OutputType.Subsumes(bodyOutputShape.Type) {
		return nullShape, nil, errors.E("type", x.Pos, "expected function body output type %s, got %s", x.OutputType, bodyOutputShape.Type)
	}
	// define output context
	outputShape := shapes.Shape{
		Type:          inputShape.Type,
		FunctionStack: functionStack,
	}
	// define action
	action := func(inputState states.State, args []functions.Action) states.State {
		bodyInputStack = inputState.Stack
		return inputState
	}
	// return
	return outputShape, action, nil
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

func (x *ConditionalExpression) Typecheck(inputShape shapes.Shape, params []*functions.Param) (shapes.Shape, functions.Action, error) {
	conditionOutputShape, conditionAction, err := x.Condition.Typecheck(inputShape, nil)
	if err != nil {
		return nullShape, nil, err
	}
	if !booleanType.Subsumes(conditionOutputShape.Type) {
		return nullShape, nil, errors.E("type", x.Pos, "condition must be boolean")
	}
	// context is the shared input context for all conditions and consequents.
	// Each condition may add to the FunctionStack. Type always stays the same.
	shape := shapes.Shape{
		Type:          inputShape.Type,
		FunctionStack: conditionOutputShape.FunctionStack,
	}
	consequentOutputShape, consequentAction, err := x.Consequent.Typecheck(shape, nil)
	if err != nil {
		return nullShape, nil, err
	}
	outputType := consequentOutputShape.Type
	elifConditionActions := make([]functions.Action, 0, len(x.ElifConditions))
	elifConsequentActions := make([]functions.Action, 0, len(x.ElifConsequents))
	for i := range x.ElifConditions {
		conditionOutputShape, elifConditionAction, err := x.ElifConditions[i].Typecheck(shape, nil)
		if err != nil {
			return nullShape, nil, err
		}
		if !booleanType.Subsumes(conditionOutputShape.Type) {
			return nullShape, nil, errors.E("type", x.Pos, "condition must be boolean")
		}
		shape.FunctionStack = conditionOutputShape.FunctionStack
		elifConditionActions = append(elifConditionActions, elifConditionAction)
		consequentOutputShape, elifConsequentAction, err := x.ElifConsequents[i].Typecheck(shape, nil)
		if err != nil {
			return nullShape, nil, err
		}
		elifConsequentActions = append(elifConsequentActions, elifConsequentAction)
		outputType = types.Disjoin(outputType, consequentOutputShape.Type)
	}
	alternativeOutputShape, alternativeAction, err := x.Alternative.Typecheck(shape, nil)
	if err != nil {
		return nullShape, nil, err
	}
	outputType = types.Disjoin(outputType, alternativeOutputShape.Type)
	action := func(inputState states.State, args []functions.Action) states.State {
		conditionState := conditionAction(inputState, nil)
		boolConditionValue, _ := conditionState.Value.(*values.BooleanValue)
		if boolConditionValue.Value {
			consequentInputState := states.State{
				Value: inputState.Value,
				Stack: conditionState.Stack,
			}
			consequentOutputState := consequentAction(consequentInputState, nil)
			return states.State{
				Value: consequentOutputState.Value,
				Stack: inputState.Stack,
			}
		}
		for i := range elifConditionActions {
			conditionState = elifConditionActions[i](inputState, nil)
			boolConditionValue, _ = conditionState.Value.(*values.BooleanValue)
			if boolConditionValue.Value {
				consequentInputState := states.State{
					Value: inputState.Value,
					Stack: conditionState.Stack,
				}
				consequentOutputState := elifConsequentActions[i](consequentInputState, nil)
				return states.State{
					Value: consequentOutputState.Value,
					Stack: inputState.Stack,
				}
			}
		}
		alternativeOutputState := alternativeAction(inputState, nil)
		return states.State{
			Value: alternativeOutputState.Value,
			Stack: inputState.Stack,
		}
	}
	outputShape := shapes.Shape{
		Type:          outputType,
		FunctionStack: inputShape.FunctionStack,
	}
	return outputShape, action, nil
}

///////////////////////////////////////////////////////////////////////////////
