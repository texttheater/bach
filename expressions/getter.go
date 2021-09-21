package expressions

import (
	"strconv"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/parameters"
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

func (x GetterExpression) Typecheck(inputShape Shape, params []*parameters.Parameter) (Shape, states.Action, *states.IDStack, error) {
	switch t := inputShape.Type.(type) {
	case types.ObjType:
		wantType := types.ObjType{
			PropTypeMap: map[string]types.Type{
				x.Name: types.AnyType{},
			},
			RestType: types.AnyType{},
		}
		if !wantType.Subsumes(inputShape.Type) {
			return Shape{}, nil, nil, errors.TypeError(
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
	case *types.NearrType:
		index, err := strconv.Atoi(x.Name)
		if err != nil {
			return Shape{}, nil, nil, errors.TypeError(
				errors.Code(errors.BadIndex),
				errors.Pos(x.Pos),
			)
		}
		var outputType types.Type
		restType := t
		if index < 0 {
			posIndex := index + 1
			revIndex := -index
			buf := make([]types.Type, revIndex)
			bufIndex := 0
			for true {
				buf[bufIndex] = restType.HeadType
				bufIndex = (bufIndex + 1) % revIndex
				if types.VoidArrType.Subsumes(restType.TailType) {
					if buf[bufIndex] == nil {
						return Shape{}, nil, nil, errors.TypeError(
							errors.Pos(x.Pos),
							errors.Code(errors.NoSuchIndex),
						)
					}
					outputType = buf[bufIndex]
					break
				}
				restType = restType.TailType.(*types.NearrType)
				posIndex += 1
			}
			index = posIndex
		} else {
			for i := 0; i < index; i++ {
				if types.VoidArrType.Subsumes(restType.TailType) {
					return Shape{}, nil, nil, errors.TypeError(
						errors.Pos(x.Pos),
						errors.Code(errors.NoSuchIndex),
					)
				}
				restType = restType.TailType.(*types.NearrType)
			}
			outputType = restType.HeadType
		}
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
		return Shape{}, nil, nil, errors.TypeError(
			errors.Code(errors.NoGetterAllowed),
			errors.Pos(x.Pos),
			errors.GotType(inputShape.Type),
		)
	}
}
