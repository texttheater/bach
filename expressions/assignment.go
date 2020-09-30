package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/parameters"
	"github.com/texttheater/bach/states"
)

type AssignmentExpression struct {
	Pos  lexer.Position
	Name string
}

func (x AssignmentExpression) Position() lexer.Position {
	return x.Pos
}

func (x AssignmentExpression) Typecheck(inputShape Shape, params []*parameters.Parameter) (Shape, states.Action, *states.IDStack, error) {
	if len(params) > 0 {
		return Shape{}, nil, nil, errors.E(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos))

	}
	variableFuncer := VariableFuncer(x, x.Name, inputShape.Type)
	outputShape := Shape{
		Type:  inputShape.Type,
		Stack: inputShape.Stack.Push(variableFuncer),
	}
	action := func(inputState states.State, args []states.Action) *states.Thunk {
		return states.ThunkFromState(states.State{
			Value: inputState.Value,
			Stack: inputState.Stack.Push(states.Variable{
				ID:     x,
				Action: states.SimpleAction(inputState.Value),
			}),
			TypeStack: inputState.TypeStack,
		})
	}
	return outputShape, action, nil, nil
}
