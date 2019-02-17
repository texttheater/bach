package builtin

import (
	"fmt"
	"os"

	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

var InitialShape = initialShape()

func initialShape() functions.Shape {
	shape := functions.Shape{types.NullType, nil}
	// math functions
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		types.NumType,
		"+",
		[]types.Type{types.NumType},
		types.NumType,
		Add,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		types.NumType,
		"-",
		[]types.Type{types.NumType},
		types.NumType,
		Subtract,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		types.NumType,
		"*",
		[]types.Type{types.NumType},
		types.NumType,
		Multiply,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		types.NumType,
		"/",
		[]types.Type{types.NumType},
		types.NumType,
		Divide,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		types.NumType,
		"%",
		[]types.Type{types.NumType},
		types.NumType,
		Modulo,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		types.NumType,
		"<",
		[]types.Type{types.NumType},
		types.BoolType,
		LessThan,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		types.NumType,
		">",
		[]types.Type{types.NumType},
		types.BoolType,
		GreaterThan,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		types.NumType,
		"==",
		[]types.Type{types.NumType},
		types.BoolType,
		Equal,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		types.NumType,
		"<=",
		[]types.Type{types.NumType},
		types.BoolType,
		LessEqual,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		types.NumType,
		">=",
		[]types.Type{types.NumType},
		types.BoolType,
		GreaterEqual,
	))
	// logic functions
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		types.AnyType,
		"true",
		nil,
		types.BoolType,
		True,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		types.AnyType,
		"false",
		nil,
		types.BoolType,
		False,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		types.BoolType,
		"and",
		[]types.Type{types.BoolType},
		types.BoolType,
		And,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		types.BoolType,
		"or",
		[]types.Type{types.BoolType},
		types.BoolType,
		Or,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		types.BoolType,
		"not",
		nil,
		types.BoolType,
		Not,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		types.BoolType,
		"==",
		[]types.Type{types.BoolType},
		types.BoolType,
		BoolEqual,
	))
	// null functions
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		types.AnyType,
		"null",
		nil,
		types.NullType,
		Null,
	))
	// sequence functions
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		types.AnyType,
		"range",
		[]types.Type{types.NumType, types.NumType},
		&types.SeqType{types.NumType},
		Range,
	))
	// conversion functions
	shape.FuncerStack = shape.FuncerStack.Push(func(gotInputType types.Type, gotName string, gotNumArgs int) ([]*functions.Parameter, types.Type, functions.Action, bool) {
		if !types.AnySeqType.Subsumes(gotInputType) {
			return nil, nil, nil, false
		}
		if gotName != "arr" {
			return nil, nil, nil, false
		}
		if gotNumArgs != 0 {
			return nil, nil, nil, false
		}
		outputType := &types.ArrType{gotInputType.ElementType()}
		action := func(inputState functions.State, args []functions.Action) functions.State {
			array := make([]values.Value, 0)
			for el := range inputState.Value.Iter() {
				array = append(array, el)
			}
			outputState := functions.State{
				Value: values.ArrValue(array),
				Stack: inputState.Stack,
			}
			return outputState
		}
		return nil, outputType, action, true
	})
	shape.FuncerStack = shape.FuncerStack.Push(func(gotInputType types.Type, gotName string, gotNumArgs int) ([]*functions.Parameter, types.Type, functions.Action, bool) {
		if gotName != "out" {
			return nil, nil, nil, false
		}
		if gotNumArgs != 0 {
			return nil, nil, nil, false
		}
		outputType := gotInputType
		action := func(inputState functions.State, args []functions.Action) functions.State {
			fmt.Println(inputState.Value)
			return inputState
		}
		return nil, outputType, action, true
	})
	shape.FuncerStack = shape.FuncerStack.Push(func(gotInputType types.Type, gotName string, gotNumArgs int) ([]*functions.Parameter, types.Type, functions.Action, bool) {
		if gotName != "err" {
			return nil, nil, nil, false
		}
		if gotNumArgs != 0 {
			return nil, nil, nil, false
		}
		outputType := gotInputType
		action := func(inputState functions.State, args []functions.Action) functions.State {
			fmt.Fprintln(os.Stderr, inputState.Value)
			return inputState
		}
		return nil, outputType, action, true
	})
	shape.FuncerStack = shape.FuncerStack.Push(func(gotInputType types.Type, gotName string, gotNumArgs int) ([]*functions.Parameter, types.Type, functions.Action, bool) {
		if gotName != "id" {
			return nil, nil, nil, false
		}
		if gotNumArgs != 0 {
			return nil, nil, nil, false
		}
		outputType := gotInputType
		action := func(inputState functions.State, args []functions.Action) functions.State {
			return inputState
		}
		return nil, outputType, action, true
	})
	return shape
}
