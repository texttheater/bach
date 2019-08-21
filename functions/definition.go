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
	funAction := func(inputState states.State, args []states.Action) states.State {
		// bind parameters to arguments by adding corresponding
		// Variable objects to the body input state.
		bodyInputStack := bodyInputStackStub
		for i, param := range x.Params {
			bodyInputStack = bodyInputStack.Push(states.Variable{
				ID:     param,
				Action: args[i],
			})
		}
		bodyInputState := states.State{
			Value: inputState.Value,
			Stack: bodyInputStack,
		}
		bodyOutputState := bodyAction(bodyInputState, nil)
		return states.State{
			Value: bodyOutputState.Value,
			Stack: inputState.Stack,
		}
	}
	// make a funcer for the defined function, add it to the function stack
	funFuncer := RegularFuncer(x.InputType, x.Name, x.Params, x.OutputType, funAction)
	functionStack := inputShape.Stack.Push(funFuncer)
	// add parameter funcers for use in the body
	bodyStack := functionStack
	for i, param := range x.Params {
		paramAction := func(inputState states.State, args []states.Action) states.State {
			stack := inputState.Stack
			for stack != nil {
				if stack.Head.ID == param {
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
			errors.Pos(x.Pos),
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
	action := func(inputState states.State, args []states.Action) states.State {
		bodyInputStackStub = inputState.Stack
		return inputState
	}
	// return
	return outputShape, action, nil
}