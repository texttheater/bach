package expressions

import (
	"strconv"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/shapes"
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

func (x GetterExpression) Typecheck(inputShape shapes.Shape, params []*params.Param) (shapes.Shape, states.Action, *states.IDStack, error) {
	switch t := inputShape.Type.(type) {
	case types.Obj:
		wantType := types.Obj{
			Props: map[string]types.Type{
				x.Name: types.Any{},
			},
			Rest: types.Any{},
		}
		if !wantType.Subsumes(inputShape.Type) {
			return shapes.Shape{}, nil, nil, errors.TypeError(
				errors.Code(errors.NoSuchProperty),
				errors.Pos(x.Pos),
				errors.WantType(wantType),
				errors.GotType(inputShape.Type),
			)
		}
		outputType := inputShape.Type.(types.Obj).Props[x.Name]
		outputShape := shapes.Shape{
			Type:  outputType,
			Stack: inputShape.Stack,
		}
		action := func(inputState states.State, args []states.Action) *states.Thunk {
			val, err := inputState.Value.(states.ObjValue)[x.Name].Eval()
			if err != nil {
				return states.ThunkFromError(err)
			}
			outputState := states.State{
				Value:     val,
				Stack:     inputState.Stack,
				TypeStack: inputState.TypeStack,
			}
			return states.ThunkFromState(outputState)
		}
		return outputShape, action, nil, nil
	case *types.Nearr:
		index, err := strconv.Atoi(x.Name)
		if err != nil || index < 0 {
			return shapes.Shape{}, nil, nil, errors.TypeError(
				errors.Code(errors.BadIndex),
				errors.Pos(x.Pos),
			)
		}
		var outputType types.Type
		restType := t
		for i := 0; i < index; i++ {
			if types.VoidArr.Subsumes(restType.Tail) {
				return shapes.Shape{}, nil, nil, errors.TypeError(
					errors.Pos(x.Pos),
					errors.Code(errors.NoSuchIndex),
				)
			}
			restType = restType.Tail.(*types.Nearr)
		}
		outputType = restType.Head
		outputShape := shapes.Shape{
			Type:  outputType,
			Stack: inputShape.Stack,
		}
		action := func(inputState states.State, args []states.Action) *states.Thunk {
			arr := inputState.Value.(*states.ArrValue)
			for i := 0; i < index; i++ {
				val, err := arr.Tail.Eval()
				if err != nil {
					return states.ThunkFromError(err)
				}
				arr = val.(*states.ArrValue)
			}
			outputState := states.State{
				Value:     arr.Head,
				Stack:     inputState.Stack,
				TypeStack: inputState.TypeStack,
			}
			return states.ThunkFromState(outputState)
		}
		return outputShape, action, nil, nil
	default:
		return shapes.Shape{}, nil, nil, errors.TypeError(
			errors.Code(errors.NoGetterAllowed),
			errors.Pos(x.Pos),
			errors.GotType(inputShape.Type),
		)
	}
}
