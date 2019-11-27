package builtin

import (
	"math"

	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initMath() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		functions.SimpleFuncer(
			types.NumType{},
			"+",
			[]types.Type{types.NumType{}},
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.NumValue(inputNum + argumentNum), nil
			},
		),
		functions.SimpleFuncer(
			types.NumType{},
			"-",
			[]types.Type{types.NumType{}},
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.NumValue(inputNum - argumentNum), nil
			},
		),
		functions.SimpleFuncer(
			types.NumType{},
			"*",
			[]types.Type{types.NumType{}},
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.NumValue(inputNum * argumentNum), nil
			},
		),
		functions.SimpleFuncer(
			types.NumType{},
			"/",
			[]types.Type{types.NumType{}},
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.NumValue(inputNum / argumentNum), nil
			},
		),
		functions.SimpleFuncer(
			types.NumType{},
			"%",
			[]types.Type{types.NumType{}},
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.NumValue(math.Mod(float64(inputNum), float64(argumentNum))), nil
			},
		),
		functions.SimpleFuncer(
			types.NumType{},
			"<",
			[]types.Type{types.NumType{}},
			types.BoolType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.BoolValue(inputNum < argumentNum), nil
			},
		),
		functions.SimpleFuncer(
			types.NumType{},
			">",
			[]types.Type{types.NumType{}},
			types.BoolType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.BoolValue(inputNum > argumentNum), nil
			},
		),
		functions.SimpleFuncer(
			types.NumType{},
			"==",
			[]types.Type{types.NumType{}},
			types.BoolType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.BoolValue(inputNum == argumentNum), nil
			},
		),
		functions.SimpleFuncer(
			types.NumType{},
			"<=",
			[]types.Type{types.NumType{}},
			types.BoolType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.BoolValue(inputNum <= argumentNum), nil
			},
		),
		functions.SimpleFuncer(
			types.NumType{},
			">=",
			[]types.Type{types.NumType{}},
			types.BoolType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputNum := inputValue.(states.NumValue)
				argumentNum := argumentValues[0].(states.NumValue)
				return states.BoolValue(inputNum >= argumentNum), nil
			},
		),
		functions.SimpleFuncer(
			&types.ArrType{types.NumType{}},
			"sum",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				arr := inputValue.(*states.ArrValue)
				sum := 0.0
				for arr != nil {
					sum += float64(arr.Head.(states.NumValue))
					var err error
					arr, err = arr.GetTail()
					if err != nil {
						return nil, err
					}
				}
				return states.NumValue(sum), nil
			},
		),
		functions.SimpleFuncer(
			&types.ArrType{types.NumType{}},
			"avg",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				arr := inputValue.(*states.ArrValue)
				sum := 0.0
				count := 0.0
				for arr != nil {
					sum += float64(arr.Head.(states.NumValue))
					count += 1.0
					var err error
					arr, err = arr.GetTail()
					if err != nil {
						return nil, err
					}
				}
				return states.NumValue(sum / count), nil
			},
		),
	})
}
