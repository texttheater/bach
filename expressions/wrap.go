package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type WrapExpression struct {
	Pos  lexer.Position
	Prop string
}

func (x WrapExpression) Position() lexer.Position {
	return x.Pos
}

func (x WrapExpression) Typecheck(inputShape shapes.Shape, params []*params.Param) (shapes.Shape, states.Action, *states.IDStack, error) {
	if len(params) > 0 {
		return shapes.Shape{}, nil, nil, errors.TypeError(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	shape := shapes.Shape{
		Type: types.Obj{
			Props: map[string]types.Type{
				x.Prop: inputShape.Type,
			},
			Rest: types.Any{},
		},
		Stack: inputShape.Stack,
	}
	action := func(inputState states.State, args []states.Action) states.State {
		wrappedValue := states.ObjFromMap(map[string]*states.Thunk{
			x.Prop: inputState.Thunk,
		})
		return inputState.Replace(states.ThunkFromValue(wrappedValue))
	}
	return shape, action, nil, nil
}
