package expressions

import (
	"fmt"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type AssignmentExpression struct {
	Pos  lexer.Position
	Name string
}

func (x AssignmentExpression) Typecheck(inputShape functions.Shape, params []*functions.Parameter) (functions.Shape, functions.Action, error) {
	if len(params) > 0 {
		return zeroShape, nil, errors.E(
			errors.Kind(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	var id interface{} = x
	varFuncer := func(gotInputType types.Type, gotName string, gotNumArgs int) ([]*functions.Parameter, types.Type, functions.Action, bool) {
		if gotName != x.Name {
			return nil, nil, nil, false
		}
		if gotNumArgs != 0 {
			return nil, nil, nil, false
		}
		varAction := func(inputState functions.State, args []functions.Action) functions.State {
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
		}
		return nil, inputShape.Type, varAction, true
	}
	outputShape := functions.Shape{inputShape.Type, inputShape.FuncerStack.Push(varFuncer)}
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
