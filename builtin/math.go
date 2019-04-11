package builtin

import (
	"math"

	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func initMath() {
	InitialShape.FuncerStack = InitialShape.FuncerStack.PushAll([]shapes.Funcer{
		shapes.SimpleFuncer(
			types.NumType(),
			"+",
			[]types.Type{types.NumType()},
			types.NumType(),
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				return values.NumValue(inputValue.Num() + argumentValues[0].Num())
			},
		),
		shapes.SimpleFuncer(
			types.NumType(),
			"-",
			[]types.Type{types.NumType()},
			types.NumType(),
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				return values.NumValue(inputValue.Num() - argumentValues[0].Num())
			},
		),
		shapes.SimpleFuncer(
			types.NumType(),
			"*",
			[]types.Type{types.NumType()},
			types.NumType(),
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				return values.NumValue(inputValue.Num() * argumentValues[0].Num())
			},
		),
		shapes.SimpleFuncer(
			types.NumType(),
			"/",
			[]types.Type{types.NumType()},
			types.NumType(),
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				return values.NumValue(inputValue.Num() / argumentValues[0].Num())
			},
		),
		shapes.SimpleFuncer(
			types.NumType(),
			"%",
			[]types.Type{types.NumType()},
			types.NumType(),
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				return values.NumValue(math.Mod(float64(inputValue.Num()), float64(argumentValues[0].Num())))
			},
		),
		shapes.SimpleFuncer(
			types.NumType(),
			"<",
			[]types.Type{types.NumType()},
			types.BoolType(),
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				return values.BoolValue(inputValue.Num() < argumentValues[0].Num())
			},
		),
		shapes.SimpleFuncer(
			types.NumType(),
			">",
			[]types.Type{types.NumType()},
			types.BoolType(),
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				return values.BoolValue(inputValue.Num() > argumentValues[0].Num())
			},
		),
		shapes.SimpleFuncer(
			types.NumType(),
			"==",
			[]types.Type{types.NumType()},
			types.BoolType(),
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				return values.BoolValue(inputValue.Num() == argumentValues[0].Num())
			},
		),
		shapes.SimpleFuncer(
			types.NumType(),
			"<=",
			[]types.Type{types.NumType()},
			types.BoolType(),
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				return values.BoolValue(inputValue.Num() <= argumentValues[0].Num())
			},
		),
		shapes.SimpleFuncer(
			types.NumType(),
			">=",
			[]types.Type{types.NumType()},
			types.BoolType(),
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				return values.BoolValue(inputValue.Num() >= argumentValues[0].Num())
			},
		),
		shapes.SimpleFuncer(
			types.SeqType(types.NumType()),
			"sum",
			nil,
			types.NumType(),
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				sum := 0.0
				for value := range inputValue.Iter() {
					sum += value.Num()
				}
				return values.NumValue(sum)
			},
		),
		shapes.SimpleFuncer(
			types.SeqType(types.NumType()),
			"average",
			nil,
			types.NumType(),
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				sum := 0.0
				count := 0.0
				for value := range inputValue.Iter() {
					sum += value.Num()
					count += 1.0
				}
				return values.NumValue(sum / count)
			},
		),
	})
}
