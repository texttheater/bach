package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type DefinitionExpression struct {
	Pos        lexer.Position
	InputType  types.Type
	Name       string
	Params     []*params.Param // these have to be pointers because we use them as IDs
	ParamNames []string
	OutputType types.Type
	Body       Expression
}

func (x DefinitionExpression) Position() lexer.Position {
	return x.Pos
}

func (x DefinitionExpression) Typecheck(inputShape shapes.Shape, params []*params.Param) (shapes.Shape, states.Action, *states.IDStack, error) {
	// make sure we got no parameters
	if len(params) > 0 {
		return shapes.Shape{}, nil, nil, errors.TypeError(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	// variables for body input stack, action (will be set at runtime)
	var bodyInputStackStub *states.VariableStack
	var bodyAction states.Action
	// make kernel for function
	funKernel := func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		// bind parameters to arguments by adding corresponding
		// Variable objects to the body input state
		bodyInputStack := bodyInputStackStub
		for i, param := range x.Params {
			arg := args[i]
			bodyInputStack = bodyInputStack.Push(states.Variable{
				ID:     param,
				Action: arg,
			})
		}
		bodyInputState := states.State{
			Value:     inputState.Value,
			Stack:     bodyInputStack,
			TypeStack: inputState.TypeStack,
		}
		return bodyAction(bodyInputState, nil)
	}
	// make a funcer for the defined function, add it to the function stack
	funFuncer := shapes.Funcer{InputType: x.InputType, Name: x.Name, Params: x.Params, OutputType: x.OutputType, Kernel: funKernel, IDs: nil}
	functionStack := inputShape.Stack.Push(funFuncer)
	// add parameter funcers for use in the body
	bodyStack := functionStack
	for i, param := range x.Params {
		id := param
		paramKernel := func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			stack := inputState.Stack
			for stack != nil {
				if stack.Head.ID == id {
					return stack.Head.Action(inputState, args)
				}
				stack = stack.Tail
			}
			panic("action not found")
		}
		paramFuncer := shapes.Funcer{InputType: param.InputType, Name: x.ParamNames[i], Params: param.Params, OutputType: param.OutputType, Kernel: paramKernel, IDs: &states.IDStack{
			Head: id,
		}}
		bodyStack = bodyStack.Push(paramFuncer)
	}
	// define body input context
	bodyInputShape := shapes.Shape{
		Type:  x.InputType,
		Stack: bodyStack,
	}
	// typecheck body (crucially, setting body action)
	bodyOutputShape, bodyAction, bodyIDs, err := x.Body.Typecheck(bodyInputShape, nil)
	if err != nil {
		return shapes.Shape{}, nil, nil, err
	}
	ids := bodyIDs
	// check body output type
	if !x.OutputType.Subsumes(bodyOutputShape.Type) {
		return shapes.Shape{}, nil, nil, errors.TypeError(
			errors.Code(errors.FunctionBodyHasWrongOutputType),
			errors.Pos(x.Pos),
			errors.WantType(x.OutputType),
			errors.GotType(bodyOutputShape.Type),
		)
	}
	// define output context
	outputShape := shapes.Shape{
		Type:  inputShape.Type,
		Stack: functionStack,
	}
	// define action (crucially, setting body input stack)
	action := func(inputState states.State, args []states.Action) *states.Thunk {
		bodyInputStackStub = inputState.Stack
		return states.ThunkFromState(inputState)

	}
	// return
	return outputShape, action, ids, nil
}
