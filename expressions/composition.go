package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/parameters"
	"github.com/texttheater/bach/states"
)

type CompositionExpression struct {
	Pos   lexer.Position
	Left  Expression
	Right Expression
}

func (x CompositionExpression) Typecheck(inputShape functions.Shape, params []*parameters.Parameter) (functions.Shape, states.Action, error) {
	if len(params) > 0 {
		return zeroShape, nil, errors.E(
			errors.Kind(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	middleShape, lAction, err := x.Left.Typecheck(inputShape, nil)
	if err != nil {
		return zeroShape, nil, err
	}
	outputShape, rAction, err := x.Right.Typecheck(middleShape, nil)
	if err != nil {
		return zeroShape, nil, err
	}
	action := func(inputState states.State, args []states.Action) states.State {
		middleState := lAction(inputState, nil)
		outputState := rAction(middleState, nil)
		return outputState
	}
	return outputShape, action, nil
}
