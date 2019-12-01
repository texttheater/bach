package functions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type RejectExpression struct {
	Pos lexer.Position
}

func (x RejectExpression) Position() lexer.Position {
	return x.Pos
}

func (x RejectExpression) Typecheck(inputShape Shape, params []*Parameter) (Shape, states.Action, error) {
	// make sure we got no parameters
	if len(params) > 0 {
		return Shape{}, nil, errors.E(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	// create output shape
	outputShape := Shape{
		Type:  types.VoidType{},
		Stack: nil,
	}
	// create action
	action := func(inputState states.State, args []states.Action) *states.Thunk {
		return states.ThunkFromError(RejectError{
			Value: inputState.Value,
		})

	}
	// return
	return outputShape, action, nil
}

type RejectError struct {
	Value states.Value
}

func (e RejectError) Error() string {
	str, err := e.Value.String()
	if err != nil {
		return err.Error()
	}
	return str + ": value rejected"
}

// ReplaceRejectError replaces a RejectError with an explainable error with
// location information about the component, and the value that was not
// handled.
func ReplaceRejectError(thunk *states.Thunk, pos lexer.Position) *states.Thunk {
	if thunk.Func == nil {
		if thunk.Result.Error == nil {
			return thunk
		}
		if rejectError, ok := thunk.Result.Error.(RejectError); ok {
			return states.ThunkFromError(errors.E(
				errors.Pos(pos),
				errors.Code(errors.UnexpectedValue),
				errors.GotValue(rejectError.Value),
			))
		}
		return thunk
	}
	return &states.Thunk{
		Func: func() *states.Thunk {
			return ReplaceRejectError(thunk.Func(), pos)
		},
	}
}
