package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type ConstantExpression struct {
	Pos   lexer.Position
	Type  types.Type
	Value values.Value
}

func (x ConstantExpression) Typecheck(inputShape functions.Shape, params []*functions.Parameter) (functions.Shape, states.Action, error) {
	if len(params) > 0 {
		return zeroShape, nil, errors.E(
			errors.Kind(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	outputShape := functions.Shape{x.Type, inputShape.FuncerStack}
	action := func(inputState states.State, args []states.Action) states.State {
		return states.State{
			Value: x.Value,
			Stack: inputState.Stack,
		}
	}
	return outputShape, action, nil
}
