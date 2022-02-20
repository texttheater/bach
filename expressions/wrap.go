package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/params"
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

func (x WrapExpression) Typecheck(inputShape Shape, params []*params.Param) (Shape, states.Action, *states.IDStack, error) {
	if len(params) > 0 {
		return Shape{}, nil, nil, errors.TypeError(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	shape := Shape{
		Type: types.Obj{
			Props: map[string]types.Type{
				x.Prop: inputShape.Type,
			},
			Rest: types.Any{},
		},
		Stack: inputShape.Stack,
	}
	action := func(inputState states.State, args []states.Action) *states.Thunk {
		wrappedValue := states.ObjValueFromMap(map[string]states.Value{
			x.Prop: inputState.Value,
		})
		return states.ThunkFromValue(wrappedValue)
	}
	return shape, action, nil, nil
}
