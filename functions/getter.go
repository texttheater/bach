package functions

import (
	"strconv"

	"github.com/alecthomas/participle/lexer"
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
	switch t := inputShape.Type.(type) {
	case types.ObjType:
		wantType := types.ObjType{
			PropTypeMap: map[string]types.Type{
				x.Name: types.AnyType{},
			},
			RestType: types.AnyType{},
		}
		if !wantType.Subsumes(inputShape.Type) {
			return Shape{}, nil, nil, states.E(
				states.Code(states.NoSuchProperty),
				states.Pos(x.Pos),
				states.WantType(wantType),
				states.GotType(inputShape.Type))

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
	case *types.NearrType:
		index, err := strconv.Atoi(x.Name)
		if err != nil {
			return Shape{}, nil, nil, states.E(
				states.Code(states.BadIndex),
				states.Pos(x.Pos))

		}
		wantType := types.AnyArrType
		for i := 0; i <= index; i++ {
			wantType = &types.NearrType{
				HeadType: types.AnyType{},
				TailType: wantType,
			}
		}
		if !wantType.Subsumes(inputShape.Type) {
			return Shape{}, nil, nil, states.E(
				states.Code(states.NoSuchIndex),
				states.WantType(wantType),
				states.GotType(inputShape.Type))

		}
		for i := 0; i < index; i++ {
			t = t.TailType.(*types.NearrType)
		}
		outputType := t.HeadType
		outputShape := Shape{
			Type:  outputType,
			Stack: inputShape.Stack,
		}
		action := func(inputState states.State, args []states.Action) *states.Thunk {
			arr := inputState.Value.(*states.ArrValue)
			for i := 0; i < index; i++ {
				result := arr.Tail.Eval()
				if result.Error != nil {
					return states.ThunkFromError(result.Error)
				}
				arr = result.Value.(*states.ArrValue)
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
		return Shape{}, nil, nil, states.E(
			states.Code(states.NoGetterAllowed),
			states.Pos(x.Pos),
			states.GotType(inputShape.Type))

	}
}
