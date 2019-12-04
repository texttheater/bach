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
				iter := states.IterFromValue(inputValue)
				sum := 0.0
				for {
					value, ok, err := iter()
					if err != nil {
						return nil, err
					}
					if !ok {
						return states.NumValue(sum), nil
					}
					sum += float64(value.(states.NumValue))
				}
			},
		),
		functions.SimpleFuncer(
			&types.ArrType{types.NumType{}},
			"avg",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				iter := states.IterFromValue(inputValue)
				sum := 0.0
				count := 0.0
				for {
					value, ok, err := iter()
					if err != nil {
						return nil, err
					}
					if !ok {
						return states.NumValue(sum / count), nil
					}
					sum += float64(value.(states.NumValue))
					count += 1.0
				}
			},
		),
		functions.SimpleFuncer(
			types.AnyType{},
			"inf",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.Inf(1)), nil
			},
		),
		functions.SimpleFuncer(
			types.AnyType{},
			"nan",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NumValue(math.NaN()), nil
			},
		),
		functions.SimpleFuncer(
			types.NumType{},
			"isFinite",
			nil,
			types.BoolType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				n := float64(inputValue.(states.NumValue))
				return states.BoolValue(!math.IsInf(n, 0)), nil
			},
		),
		functions.SimpleFuncer(
			types.NumType{},
			"isNaN",
			nil,
			types.BoolType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				n := float64(inputValue.(states.NumValue))
				return states.BoolValue(math.IsNaN(n)), nil
			},
		),
	})
}
