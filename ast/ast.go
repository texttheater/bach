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
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

///////////////////////////////////////////////////////////////////////////////

var nullShape = functions.Shape{}

var booleanType = types.BooleanType{}

///////////////////////////////////////////////////////////////////////////////

type Expression interface {
	Typecheck(inputShape functions.Shape, params []*functions.Parameter) (outputShape functions.Shape, action functions.Action, err error)
}

///////////////////////////////////////////////////////////////////////////////

type ConstantExpression struct {
	Pos   lexer.Position
	Type  types.Type
	Value values.Value
}

func (x *ConstantExpression) Typecheck(inputShape functions.Shape, params []*functions.Parameter) (functions.Shape, functions.Action, error) {
	if len(params) > 0 {
		return nullShape, nil, errors.E("type", x.Pos, "number expression does not take parameters")
	}
	outputShape := functions.Shape{x.Type, inputShape.FunctionStack}
	action := func(inputState functions.State, args []functions.Action) functions.State {
		return functions.State{
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

func (x *CompositionExpression) Typecheck(inputShape functions.Shape, params []*functions.Parameter) (functions.Shape, functions.Action, error) {
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
	action := func(inputState functions.State, args []functions.Action) functions.State {
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

func (x *CallExpression) Typecheck(inputShape functions.Shape, params []*functions.Parameter) (functions.Shape, functions.Action, error) {
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
			argInputShape := functions.Shape{param.InputType, inputShape.FunctionStack}
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
		return functions.Shape{function.OutputType, inputShape.FunctionStack}, action, nil
	}
}

///////////////////////////////////////////////////////////////////////////////

type AssignmentExpression struct {
	Pos  lexer.Position
	Name string
}

func (x *AssignmentExpression) Typecheck(inputShape functions.Shape, params []*functions.Parameter) (functions.Shape, functions.Action, error) {
	if len(params) > 0 {
		return nullShape, nil, errors.E("type", x.Pos, "assignment expression does not take parameters")
	}
	var id interface{} = x
	outputShape := functions.Shape{inputShape.Type, inputShape.FunctionStack.Push(functions.Function{
		InputType:  &types.AnyType{},
		Name:       x.Name,
		Params:     nil,
		OutputType: inputShape.Type,
		Action: func(inputState functions.State, args []functions.Action) functions.State {
			stack := inputState.Stack
			for stack != nil {
				if stack.Head.ID == id {
					return functions.State{
						Value: stack.Head.Action(functions.InitialState, nil).Value,
						Stack: inputState.Stack,
					}
				}
				stack = stack.Tail
			}
			panic(fmt.Sprintf("variable %s not found", x.Name))
		},
	})}
	action := func(inputState functions.State, args []functions.Action) functions.State {
		return functions.State{
			Value: inputState.Value,
			Stack: inputState.Stack.Push(functions.Variable{
				ID: id,
				Action: func(i functions.State, as []functions.Action) functions.State {
					return functions.State{
						Value: inputState.Value,
						Stack: i.Stack,
					}
				},
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
	Params     []*functions.Parameter
	OutputType types.Type
	Body       Expression
}

func (x *DefinitionExpression) Typecheck(inputShape functions.Shape, params []*functions.Parameter) (functions.Shape, functions.Action, error) {
	// make sure we got no parameters
	if len(params) > 0 {
		return nullShape, nil, errors.E("type", x.Pos, "definition expression does not take parameters")
	}
	// variables for body input stack, action (will be set at runtime)
	var bodyInputStackStub *functions.VariableStack
	var bodyAction functions.Action
	// add the function defined here to the function stack
	functionStack := inputShape.FunctionStack.Push(functions.Function{
		InputType:  x.InputType,
		Name:       x.Name,
		Params:     x.Params,
		OutputType: x.OutputType,
		Action: func(inputState functions.State, args []functions.Action) functions.State {
			// Bind parameters to arguments by adding corresponding
			// Variable objects to the body input state.
			bodyInputStack := bodyInputStackStub
			for i, param := range x.Params {
				var id interface{} = param
				bodyInputStack = bodyInputStack.Push(functions.Variable{
					ID:     id,
					Action: args[i],
				})
			}
			bodyInputState := functions.State{
				Value: inputState.Value,
				Stack: bodyInputStack,
			}
			bodyOutputState := bodyAction(bodyInputState, nil)
			return functions.State{
				Value: bodyOutputState.Value,
				Stack: inputState.Stack,
			}
		},
	})
	// add parameter functions for use in the body
	bodyFunctionStack := functionStack
	for _, param := range x.Params {
		var id interface{} = param
		bodyFunctionStack = bodyFunctionStack.Push(functions.Function{
			InputType:  param.InputType,
			Name:       param.Name,
			Params:     param.Params,
			OutputType: param.OutputType,
			Action: func(inputState functions.State, args []functions.Action) functions.State {
				stack := inputState.Stack
				for stack != nil {
					if stack.Head.ID == id {
						return stack.Head.Action(inputState, args)
					}
					stack = stack.Tail
				}
				panic("action not found")
			},
		})
	}
	// define body input context
	bodyInputShape := functions.Shape{
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
	outputShape := functions.Shape{
		Type:          inputShape.Type,
		FunctionStack: functionStack,
	}
	// define action (crucially, setting body input stack)
	action := func(inputState functions.State, args []functions.Action) functions.State {
		bodyInputStackStub = inputState.Stack
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

func (x *ConditionalExpression) Typecheck(inputShape functions.Shape, params []*functions.Parameter) (functions.Shape, functions.Action, error) {
	conditionOutputShape, conditionAction, err := x.Condition.Typecheck(inputShape, nil)
	if err != nil {
		return nullShape, nil, err
	}
	if !booleanType.Subsumes(conditionOutputShape.Type) {
		return nullShape, nil, errors.E("type", x.Pos, "condition must be boolean")
	}
	// context is the shared input context for all conditions and consequents.
	// Each condition may add to the FunctionStack. Type always stays the same.
	shape := functions.Shape{
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
	action := func(inputState functions.State, args []functions.Action) functions.State {
		conditionState := conditionAction(inputState, nil)
		boolConditionValue, _ := conditionState.Value.(*values.BooleanValue)
		if boolConditionValue.Value {
			consequentInputState := functions.State{
				Value: inputState.Value,
				Stack: conditionState.Stack,
			}
			consequentOutputState := consequentAction(consequentInputState, nil)
			return functions.State{
				Value: consequentOutputState.Value,
				Stack: inputState.Stack,
			}
		}
		for i := range elifConditionActions {
			conditionState = elifConditionActions[i](inputState, nil)
			boolConditionValue, _ = conditionState.Value.(*values.BooleanValue)
			if boolConditionValue.Value {
				consequentInputState := functions.State{
					Value: inputState.Value,
					Stack: conditionState.Stack,
				}
				consequentOutputState := elifConsequentActions[i](consequentInputState, nil)
				return functions.State{
					Value: consequentOutputState.Value,
					Stack: inputState.Stack,
				}
			}
		}
		alternativeOutputState := alternativeAction(inputState, nil)
		return functions.State{
			Value: alternativeOutputState.Value,
			Stack: inputState.Stack,
		}
	}
	outputShape := functions.Shape{
		Type:          outputType,
		FunctionStack: inputShape.FunctionStack,
	}
	return outputShape, action, nil
}

///////////////////////////////////////////////////////////////////////////////
