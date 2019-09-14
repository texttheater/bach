package functions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/values"
)

type AssignmentExpression struct {
	Pos  lexer.Position
	Name string
}

func (x AssignmentExpression) Position() lexer.Position {
	return x.Pos
}

func (x AssignmentExpression) Typecheck(inputShape Shape, params []*Parameter) (Shape, states.Action, error) {
	if len(params) > 0 {
		return Shape{}, nil, errors.E(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	variableFuncer := VariableFuncer(x, x.Name, inputShape.Type)
	outputShape := Shape{inputShape.Type, inputShape.Stack.Push(variableFuncer)}
	action := func(inputState states.State, args []states.Action) states.State {
		return states.State{
			Value: inputState.Value,
			Stack: inputState.Stack.Push(states.Variable{
				ID:     x,
				Action: states.SimpleAction(inputState.Value),
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
