package expressions

import (
	"fmt"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/parameters"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type AssignmentExpression struct {
	Pos  lexer.Position
	Name string
}

func (x AssignmentExpression) Typecheck(inputShape shapes.Shape, params []*parameters.Parameter) (shapes.Shape, states.Action, error) {
	if len(params) > 0 {
		return zeroShape, nil, errors.E(
			errors.Kind(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	var id interface{} = x
	varFuncer := func(gotInputType types.Type, gotName string, gotNumArgs int) ([]*parameters.Parameter, types.Type, states.Action, bool) {
		if gotName != x.Name {
			return nil, nil, nil, false
		}
		if gotNumArgs != 0 {
			return nil, nil, nil, false
		}
		varAction := func(inputState states.State, args []states.Action) states.State {
			stack := inputState.Stack
			for stack != nil {
				if stack.Head.ID == id {
					return states.State{
						Value: stack.Head.Action(states.InitialState, nil).Value,
						Stack: inputState.Stack,
					}
				}
				stack = stack.Tail
			}
			panic(fmt.Sprintf("variable %s not found", x.Name))
		}
		return nil, inputShape.Type, varAction, true
	}
	outputShape := shapes.Shape{inputShape.Type, inputShape.FuncerStack.Push(varFuncer)}
	action := func(inputState states.State, args []states.Action) states.State {
		return states.State{
			Value: inputState.Value,
			Stack: inputState.Stack.Push(states.Variable{
				ID: id,
				Action: func(i states.State, as []states.Action) states.State {
					return states.State{
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
