package functions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type DefinitionExpression struct {
	Pos        lexer.Position
	InputType  types.Type
	Name       string
	Params     []*Parameter
	ParamNames []string
	OutputType types.Type
	Body       Expression
}

func (x DefinitionExpression) Position() lexer.Position {
	return x.Pos
}

func (x DefinitionExpression) Typecheck(inputShape Shape, params []*Parameter) (Shape, states.Action, error) {
	// make sure we got no parameters
	if len(params) > 0 {
		return Shape{}, nil, errors.E(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	// variables for body input stack, action (will be set at runtime)
	var bodyInputStackStub *states.VariableStack
	var bodyAction states.Action
	// make action for function
	funAction := func(inputState states.State, args []states.Action) (states.State, bool, error) {
		// bind parameters to arguments by adding corresponding
		// Variable objects to the body input state
		bodyInputStack := bodyInputStackStub
		for i, param := range x.Params {
			arg := args[i]
			bodyInputStack = bodyInputStack.Push(states.Variable{
				ID: param,
				Action: func(argInputState states.State, argArgs []states.Action) (states.State, bool, error) {
					outputState, _, err := arg(argInputState, argArgs)
					if err != nil {
						return states.State{}, false, err
					}
					outputState.Stack = argInputState.Stack
					outputState.TypeStack = argInputState.TypeStack // TODO right?
					return outputState, false, nil
				},
			})
		}
		bodyInputState := states.State{
			Value:     inputState.Value,
			Stack:     bodyInputStack,
			TypeStack: inputState.TypeStack,
		}
		//return func() (states.State, bool, error, states.Thunk) {
		//	return bodyAction(bodyInputState, nil)
		//}
		bodyOutputState, _, err := bodyAction(bodyInputState, nil)
		if err != nil {
			return states.State{}, false, err
		}
		return states.State{
			Value:     bodyOutputState.Value,
			Stack:     inputState.Stack,
			TypeStack: inputState.TypeStack,
		}, false, nil
	}
	// make a funcer for the defined function, add it to the function stack
	funFuncer := RegularFuncer(x.InputType, x.Name, x.Params, x.OutputType, funAction)
	functionStack := inputShape.Stack.Push(funFuncer)
	// add parameter funcers for use in the body
	bodyStack := functionStack
	for i, param := range x.Params {
		id := param
		paramAction := func(inputState states.State, args []states.Action) (states.State, bool, error) {
			stack := inputState.Stack
			for stack != nil {
				if stack.Head.ID == id {
					return stack.Head.Action(inputState, args)
				}
				stack = stack.Tail
			}
			panic("action not found")
		}
		paramFuncer := RegularFuncer(param.InputType, x.ParamNames[i], param.Params, param.OutputType, paramAction)
		bodyStack = bodyStack.Push(paramFuncer)
	}
	// define body input context
	bodyInputShape := Shape{
		Type:  x.InputType,
		Stack: bodyStack,
	}
	// typecheck body (crucially, setting body action)
	bodyOutputShape, bodyAction, err := x.Body.Typecheck(bodyInputShape, nil)
	if err != nil {
		return Shape{}, nil, err
	}
	// check body output type
	if !x.OutputType.Subsumes(bodyOutputShape.Type) {
		return Shape{}, nil, errors.E(
			errors.Code(errors.FunctionBodyHasWrongOutputType),
			errors.Pos(x.Body.Position()),
			errors.WantType(x.OutputType),
			errors.GotType(bodyOutputShape.Type),
		)
	}
	// define output context
	outputShape := Shape{
		Type:  inputShape.Type,
		Stack: functionStack,
	}
	// define action (crucially, setting body input stack)
	action := func(inputState states.State, args []states.Action) (states.State, bool, error) {
		bodyInputStackStub = inputState.Stack
		return inputState, false, nil
	}
	// return
	return outputShape, action, nil
}
