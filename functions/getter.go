package functions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type GetterExpression struct {
	Pos  lexer.Position
	Name string
}

func (x GetterExpression) Position() lexer.Position {
	return x.Pos
}

func (x GetterExpression) Typecheck(inputShape Shape, params []*Parameter) (Shape, states.Action, *states.IDStack, error) {
	wantType := types.ObjType{
		PropTypeMap: map[string]types.Type{
			x.Name: types.AnyType{},
		},
		RestType: types.AnyType{},
	}
	if !wantType.Subsumes(inputShape.Type) {
		return Shape{}, nil, nil, errors.E(
			errors.Code(errors.NoSuchProperty),
			errors.Pos(x.Pos),
			errors.WantType(wantType),
			errors.GotType(inputShape.Type),
		)
	}
	outputType := inputShape.Type.(types.ObjType).PropTypeMap[x.Name]
	outputShape := Shape{
		Type:  outputType,
		Stack: inputShape.Stack,
	}
	action := func(inputState states.State, args []states.Action) *states.Thunk {
		result := inputState.Value.(states.ObjValue)[x.Name].Eval()
		if result.Error != nil {
			return states.ThunkFromError(result.Error)
		}
		outputState := states.State{
			Value:     result.Value,
			Stack:     inputState.Stack,
			TypeStack: inputState.TypeStack,
		}
		return states.ThunkFromState(outputState)
	}
	return outputShape, action, nil, nil
}
