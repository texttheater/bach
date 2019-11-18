package builtin

import (
	"math"

	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func initMath() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		functions.SimpleFuncer(
			types.NumType{},
			"+",
			[]types.Type{types.NumType{}},
			types.NumType{},
			func(inputValue values.Value, argumentValues []values.Value) (values.Value, error) {
				inputNum := inputValue.(values.NumValue)
				argumentNum := argumentValues[0].(values.NumValue)
				return values.NumValue(inputNum + argumentNum), nil
			},
		),
		functions.SimpleFuncer(
			types.NumType{},
			"-",
			[]types.Type{types.NumType{}},
			types.NumType{},
			func(inputValue values.Value, argumentValues []values.Value) (values.Value, error) {
				inputNum := inputValue.(values.NumValue)
				argumentNum := argumentValues[0].(values.NumValue)
				return values.NumValue(inputNum - argumentNum), nil
			},
		),
		functions.SimpleFuncer(
			types.NumType{},
			"*",
			[]types.Type{types.NumType{}},
			types.NumType{},
			func(inputValue values.Value, argumentValues []values.Value) (values.Value, error) {
				inputNum := inputValue.(values.NumValue)
				argumentNum := argumentValues[0].(values.NumValue)
				return values.NumValue(inputNum * argumentNum), nil
			},
		),
		functions.SimpleFuncer(
			types.NumType{},
			"/",
			[]types.Type{types.NumType{}},
			types.NumType{},
			func(inputValue values.Value, argumentValues []values.Value) (values.Value, error) {
				inputNum := inputValue.(values.NumValue)
				argumentNum := argumentValues[0].(values.NumValue)
				return values.NumValue(inputNum / argumentNum), nil
			},
		),
		functions.SimpleFuncer(
			types.NumType{},
			"%",
			[]types.Type{types.NumType{}},
			types.NumType{},
			func(inputValue values.Value, argumentValues []values.Value) (values.Value, error) {
				inputNum := inputValue.(values.NumValue)
				argumentNum := argumentValues[0].(values.NumValue)
				return values.NumValue(math.Mod(float64(inputNum), float64(argumentNum))), nil
			},
		),
		functions.SimpleFuncer(
			types.NumType{},
			"<",
			[]types.Type{types.NumType{}},
			types.BoolType{},
			func(inputValue values.Value, argumentValues []values.Value) (values.Value, error) {
				inputNum := inputValue.(values.NumValue)
				argumentNum := argumentValues[0].(values.NumValue)
				return values.BoolValue(inputNum < argumentNum), nil
			},
		),
		functions.SimpleFuncer(
			types.NumType{},
			">",
			[]types.Type{types.NumType{}},
			types.BoolType{},
			func(inputValue values.Value, argumentValues []values.Value) (values.Value, error) {
				inputNum := inputValue.(values.NumValue)
				argumentNum := argumentValues[0].(values.NumValue)
				return values.BoolValue(inputNum > argumentNum), nil
			},
		),
		functions.SimpleFuncer(
			types.NumType{},
			"==",
			[]types.Type{types.NumType{}},
			types.BoolType{},
			func(inputValue values.Value, argumentValues []values.Value) (values.Value, error) {
				inputNum := inputValue.(values.NumValue)
				argumentNum := argumentValues[0].(values.NumValue)
				return values.BoolValue(inputNum == argumentNum), nil
			},
		),
		functions.SimpleFuncer(
			types.NumType{},
			"<=",
			[]types.Type{types.NumType{}},
			types.BoolType{},
			func(inputValue values.Value, argumentValues []values.Value) (values.Value, error) {
				inputNum := inputValue.(values.NumValue)
				argumentNum := argumentValues[0].(values.NumValue)
				return values.BoolValue(inputNum <= argumentNum), nil
			},
		),
		functions.SimpleFuncer(
			types.NumType{},
			">=",
			[]types.Type{types.NumType{}},
			types.BoolType{},
			func(inputValue values.Value, argumentValues []values.Value) (values.Value, error) {
				inputNum := inputValue.(values.NumValue)
				argumentNum := argumentValues[0].(values.NumValue)
				return values.BoolValue(inputNum >= argumentNum), nil
			},
		),
		functions.SimpleFuncer(
			&types.ArrType{types.NumType{}},
			"sum",
			nil,
			types.NumType{},
			func(inputValue values.Value, argumentValues []values.Value) (values.Value, error) {
				arr := inputValue.(*values.ArrValue)
				sum := 0.0
				for {
					err := arr.Eval()
					if err != nil {
						return nil, err
					}
					if arr.Head == nil {
						break
					}
					sum += float64(arr.Head.(values.NumValue))
					arr = arr.Tail
				}
				return values.NumValue(sum), nil
			},
		),
		functions.SimpleFuncer(
			&types.ArrType{types.NumType{}},
			"avg",
			nil,
			types.NumType{},
			func(inputValue values.Value, argumentValues []values.Value) (values.Value, error) {
				arr := inputValue.(*values.ArrValue)
				sum := 0.0
				count := 0.0
				for {
					err := arr.Eval()
					if err != nil {
						return nil, err
					}
					if arr.Head == nil {
						break
					}
					sum += float64(arr.Head.(values.NumValue))
					count += 1.0
				}
				return values.NumValue(sum / count), nil
			},
		),
	})
}
