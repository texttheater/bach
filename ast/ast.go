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
	"github.com/texttheater/bach/parameters"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

///////////////////////////////////////////////////////////////////////////////

var nullShape = shapes.Shape{}

var nullAction = states.Action{}

var booleanType = types.BooleanType{}

var debug bool = true

///////////////////////////////////////////////////////////////////////////////

type Expression interface {
	Typecheck(inputShape shapes.Shape, params []*parameters.Param) (outputShape shapes.Shape, action states.Action, err error)
}

///////////////////////////////////////////////////////////////////////////////

type ConstantExpression struct {
	Pos   lexer.Position
	Type  types.Type
	Value values.Value
}

func (x *ConstantExpression) Typecheck(inputShape shapes.Shape, params []*parameters.Param) (shapes.Shape, states.Action, error) {
	if len(params) > 0 {
		return nullShape, nullAction, errors.E("type", x.Pos, "number expression does not take parameters")
	}
	outputShape := shapes.Shape{x.Type, inputShape.FunctionStack}
	action := states.Action{
		Name: "",
		Execute: func(inputState states.State, args []states.Action) states.State {
			return states.State{
				Value:       x.Value,
				ActionStack: inputState.ActionStack,
			}
		},
	}
	return outputShape, action, nil
}

///////////////////////////////////////////////////////////////////////////////

type CompositionExpression struct {
	Pos   lexer.Position
	Left  Expression
	Right Expression
}

func (x *CompositionExpression) Typecheck(inputShape shapes.Shape, params []*parameters.Param) (shapes.Shape, states.Action, error) {
	if len(params) > 0 {
		return nullShape, nullAction, errors.E("type", x.Pos, "composition expression does not take parameters")
	}
	middleShape, lAction, err := x.Left.Typecheck(inputShape, nil)
	if err != nil {
		return nullShape, nullAction, err
	}
	outputShape, rAction, err := x.Right.Typecheck(middleShape, nil)
	if err != nil {
		return nullShape, nullAction, err
	}
	action := states.Action{
		Name: "",
		Execute: func(inputState states.State, args []states.Action) states.State {
			middleState := lAction.Execute(inputState, nil)
			outputState := rAction.Execute(middleState, nil)
			return outputState
		},
	}
	return outputShape, action, nil
}

///////////////////////////////////////////////////////////////////////////////

type CallExpression struct {
	Pos  lexer.Position
	Name string
	Args []Expression
}

func (x *CallExpression) Typecheck(inputShape shapes.Shape, params []*parameters.Param) (shapes.Shape, states.Action, error) {
	// Go down the function stack and find the function invoked by this
	// call
	stack := inputShape.FunctionStack
Entries:
	for {
		// Reached bottom of stack without finding a matching function
		if stack == nil {
			return nullShape, nullAction, errors.E("type", x.Pos, "no such function")
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

func (x *AssignmentExpression) Typecheck(inputShape shapes.Shape, params []*parameters.Param) (shapes.Shape, states.Action, error) {
	if len(params) > 0 {
		return nullShape, nullAction, errors.E("type", x.Pos, "assignment expression does not take parameters")
	}
	outputShape := shapes.Shape{inputShape.Type, inputShape.FunctionStack.Push(functions.Function{
		InputType:  &types.AnyType{},
		Name:       x.Name,
		Params:     nil,
		OutputType: inputShape.Type,
		Action: states.Action{
			Name: x.Name,
			Execute: func(inputState states.State, args []states.Action) states.State {
				stack := inputState.ActionStack
				for stack != nil {
					if stack.Head.Name == x.Name {
						return stack.Head.Execute(inputState, args)
					}
					stack = stack.Tail
				}
				panic("unknown variable")
			},
		},
	})}
	action := states.Action{
		Name: "",
		Execute: func(assInputState states.State, args []states.Action) states.State {
			return states.State{
				Value: assInputState.Value,
				ActionStack: assInputState.ActionStack.Push(states.Action{
					Name: x.Name,
					Execute: func(varInputState states.State, args []states.Action) states.State {
						return states.State{
							Value:       assInputState.Value,
							ActionStack: varInputState.ActionStack,
						}
					},
				}),
			}
		},
	}
	return outputShape, action, nil
}

///////////////////////////////////////////////////////////////////////////////

type DefinitionExpression struct {
	Pos        lexer.Position
	InputType  types.Type
	Name       string
	Params     []*parameters.Param
	OutputType types.Type
	Body       Expression
}

func (x *DefinitionExpression) Typecheck(inputShape shapes.Shape, params []*parameters.Param) (shapes.Shape, states.Action, error) {
	// make sure we got no parameters
	if len(params) > 0 {
		return nullShape, nullAction, errors.E("type", x.Pos, "definition expression does not take parameters")
	}
	// dummy body action (will be set later)
	bodyAction := nullAction
	// add the function defined here to the function stack
	functionStack := inputShape.FunctionStack.Push(functions.Function{
		InputType:  x.InputType,
		Name:       x.Name,
		Params:     x.Params,
		OutputType: x.OutputType,
		Action: states.Action{
			Name: x.Name,
			Execute: func(inputState states.State, args []states.Action) states.State {
				actionStack := inputState.ActionStack
				for i, param := range x.Params {
					actionStack = actionStack.Push(states.Action{
						Name:    param.Name,
						Execute: args[i].Execute,
					})
				}
				return bodyAction.Execute(states.State{
					Value:       inputState.Value,
					ActionStack: actionStack,
				}, nil)
			},
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
			Action: states.Action{
				Name: param.Name,
				Execute: func(inputState states.State, args []states.Action) states.State {
					stack := inputState.ActionStack
					for stack != nil {
						if stack.Head.Name == param.Name {
							return stack.Head.Execute(inputState, args)
						}
						stack = stack.Tail
					}
					panic("param action not found")
				},
			},
		})
	}
	// define body input shape
	bodyInputShape := shapes.Shape{
		Type:          x.InputType,
		FunctionStack: bodyFunctionStack,
	}
	// typecheck body (crucially, setting body action)
	bodyOutputShape, bodyAction, err := x.Body.Typecheck(bodyInputShape, nil)
	if err != nil {
		return nullShape, nullAction, err
	}
	// check body output type
	if !x.OutputType.Subsumes(bodyOutputShape.Type) {
		return nullShape, nullAction, errors.E("type", x.Pos, "expected function body output type %s, got %s", x.OutputType, bodyOutputShape.Type)
	}
	// define output shape
	outputShape := shapes.Shape{
		Type:          inputShape.Type,
		FunctionStack: functionStack,
	}
	// define action (simple identity)
	action := states.Action{
		Name: "",
		Execute: func(inputState states.State, args []states.Action) states.State {
			return inputState
		},
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

func (x *ConditionalExpression) Typecheck(inputShape shapes.Shape, params []*parameters.Param) (shapes.Shape, states.Action, error) {
	conditionOutputShape, conditionAction, err := x.Condition.Typecheck(inputShape, nil)
	if err != nil {
		return nullShape, nullAction, err
	}
	if !booleanType.Subsumes(conditionOutputShape.Type) {
		return nullShape, nullAction, errors.E("type", x.Pos, "condition must be boolean")
	}
	// shape is the shared input shape for all conditions and consequents.
	// Each condition may add to the FunctionStack. Type always stays the same.
	shape := shapes.Shape{
		Type:          inputShape.Type,
		FunctionStack: conditionOutputShape.FunctionStack,
	}
	consequentOutputShape, consequentAction, err := x.Consequent.Typecheck(shape, nil)
	if err != nil {
		return nullShape, nullAction, err
	}
	outputType := consequentOutputShape.Type
	elifConditionActions := make([]states.Action, 0, len(x.ElifConditions))
	elifConsequentActions := make([]states.Action, 0, len(x.ElifConsequents))
	for i := range x.ElifConditions {
		conditionOutputShape, elifConditionAction, err := x.ElifConditions[i].Typecheck(shape, nil)
		if err != nil {
			return nullShape, nullAction, err
		}
		if !booleanType.Subsumes(conditionOutputShape.Type) {
			return nullShape, nullAction, errors.E("type", x.Pos, "condition must be boolean")
		}
		shape.FunctionStack = conditionOutputShape.FunctionStack
		elifConditionActions = append(elifConditionActions, elifConditionAction)
		consequentOutputShape, elifConsequentAction, err := x.ElifConsequents[i].Typecheck(shape, nil)
		if err != nil {
			return nullShape, nullAction, err
		}
		elifConsequentActions = append(elifConsequentActions, elifConsequentAction)
		outputType = types.Disjoin(outputType, consequentOutputShape.Type)
	}
	alternativeOutputShape, alternativeAction, err := x.Alternative.Typecheck(shape, nil)
	if err != nil {
		return nullShape, nullAction, err
	}
	outputType = types.Disjoin(outputType, alternativeOutputShape.Type)
	action := states.Action{
		Name: "",
		Execute: func(inputState states.State, args []states.Action) states.State {
			state := inputState
			conditionOutputState := conditionAction.Execute(inputState, nil)
			state.ActionStack = conditionOutputState.ActionStack
			conditionValue := conditionOutputState.Value
			boolConditionValue, _ := conditionValue.(*values.BooleanValue)
			if boolConditionValue.Value {
				consequentOutputState := consequentAction.Execute(state, nil)
				return states.State{
					Value:       consequentOutputState.Value,
					ActionStack: inputState.ActionStack,
				}
			}
			for i := range elifConditionActions {
				elifConditionOutputState := elifConditionActions[i].Execute(state, nil)
				state.ActionStack = elifConditionOutputState.ActionStack
				conditionValue = elifConditionOutputState.Value
				boolConditionValue, _ = conditionValue.(*values.BooleanValue)
				if boolConditionValue.Value {
					elifConsequentOutputState := elifConsequentActions[i].Execute(state, nil)
					return states.State{
						Value:       elifConsequentOutputState.Value,
						ActionStack: inputState.ActionStack,
					}
				}
			}
			alternativeOutputState := alternativeAction.Execute(state, nil)
			return states.State{
				Value:       alternativeOutputState.Value,
				ActionStack: inputState.ActionStack,
			}
		},
	}
	outputShape := shapes.Shape{
		Type:          outputType,
		FunctionStack: inputShape.FunctionStack,
	}
	return outputShape, action, nil
}

///////////////////////////////////////////////////////////////////////////////
