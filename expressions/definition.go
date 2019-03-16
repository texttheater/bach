package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
)

type DefinitionExpression struct {
	Pos        lexer.Position
	InputType  types.Type
	Name       string
	Params     []*functions.Parameter
	OutputType types.Type
	Body       Expression
}

func (x DefinitionExpression) Typecheck(inputShape functions.Shape, params []*functions.Parameter) (functions.Shape, functions.Action, error) {
	// make sure we got no parameters
	if len(params) > 0 {
		return zeroShape, nil, errors.E(
			errors.Kind(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	// variables for body input stack, action (will be set at runtime)
	var bodyInputStackStub *functions.VariableStack
	var bodyAction functions.Action
	// make a funcer for the defined function, add it to the function stack
	funFuncer := func(gotInputType types.Type, gotName string, gotNumArgs int) ([]*functions.Parameter, types.Type, functions.Action, bool) {
		if !x.InputType.Subsumes(gotInputType) {
			return nil, nil, nil, false
		}
		if gotName != x.Name {
			return nil, nil, nil, false
		}
		if gotNumArgs != len(x.Params) {
			return nil, nil, nil, false
		}
		funAction := func(inputState functions.State, args []functions.Action) functions.State {
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
		}
		return x.Params, x.OutputType, funAction, true
	}
	functionStack := inputShape.FuncerStack.Push(funFuncer)
	// add parameter funcers for use in the body
	bodyFuncerStack := functionStack
	for _, param := range x.Params {
		var id interface{} = param
		paramFuncer := func(gotInputType types.Type, gotName string, gotNumArgs int) ([]*functions.Parameter, types.Type, functions.Action, bool) {
			if !param.InputType.Subsumes(gotInputType) {
				return nil, nil, nil, false
			}
			if gotName != param.Name {
				return nil, nil, nil, false
			}
			if gotNumArgs != len(param.Params) {
				return nil, nil, nil, false
			}
			paramAction := func(inputState functions.State, args []functions.Action) functions.State {
				stack := inputState.Stack
				for stack != nil {
					if stack.Head.ID == id {
						return stack.Head.Action(inputState, args)
					}
					stack = stack.Tail
				}
				panic("action not found")
			}
			return param.Params, param.OutputType, paramAction, true
		}

		bodyFuncerStack = bodyFuncerStack.Push(paramFuncer)
	}
	// define body input context
	bodyInputShape := functions.Shape{
		Type:        x.InputType,
		FuncerStack: bodyFuncerStack,
	}
	// typecheck body (crucially, setting body action)
	bodyOutputShape, bodyAction, err := x.Body.Typecheck(bodyInputShape, nil)
	if err != nil {
		return zeroShape, nil, err
	}
	// check body output type
	if !x.OutputType.Subsumes(bodyOutputShape.Type) {
		return zeroShape, nil, errors.E(
			errors.Kind(errors.FunctionBodyHasWrongOutputType),
			errors.Pos(x.Pos),
			errors.WantType(x.OutputType),
			errors.GotType(bodyOutputShape.Type),
		)
	}
	// define output context
	outputShape := functions.Shape{
		Type:        inputShape.Type,
		FuncerStack: functionStack,
	}
	// define action (crucially, setting body input stack)
	action := func(inputState functions.State, args []functions.Action) functions.State {
		bodyInputStackStub = inputState.Stack
		return inputState
	}
	// return
	return outputShape, action, nil
}
