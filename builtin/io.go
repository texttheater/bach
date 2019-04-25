package builtin

import (
	"bufio"
	"fmt"
	"os"

	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func initIO() {
	InitialShape.FuncerStack = InitialShape.FuncerStack.PushAll([]shapes.Funcer{
		shapes.SimpleFuncer(
			types.AnyType{},
			"in",
			nil,
			types.ReaderType{},
			func(inputValue values.Value, argValues []values.Value) values.Value {
				return values.ReaderValue{
					os.Stdin,
				}
			},
		),
		shapes.SimpleFuncer(
			types.ReaderType{},
			"lines",
			nil,
			&types.SeqType{types.StrType{}},
			func(inputValue values.Value, argValues []values.Value) values.Value {
				reader, _ := inputValue.(values.ReaderValue)
				lines := make(chan values.Value)
				scanner := bufio.NewScanner(reader.Reader)
				go func() {
					for scanner.Scan() {
						lines <- values.StrValue(scanner.Text())
					}
					close(lines)
				}()
				return &values.SeqValue{types.StrType{}, lines}
			},
		),
		func(gotInputType types.Type, gotName string, gotNumArgs int) ([]*shapes.Parameter, types.Type, states.Action, bool) {
			if gotName != "out" {
				return nil, nil, nil, false
			}
			if gotNumArgs != 0 {
				return nil, nil, nil, false
			}
			outputType := gotInputType
			action := func(inputState states.State, args []states.Action) states.State {
				fmt.Println(inputState.Value.Out())
				return inputState
			}
			return nil, outputType, action, true
		},
		func(gotInputType types.Type, gotName string, gotNumArgs int) ([]*shapes.Parameter, types.Type, states.Action, bool) {
			if gotName != "err" {
				return nil, nil, nil, false
			}
			if gotNumArgs != 0 {
				return nil, nil, nil, false
			}
			outputType := gotInputType
			action := func(inputState states.State, args []states.Action) states.State {
				fmt.Fprintln(os.Stderr, inputState.Value)
				return inputState
			}
			return nil, outputType, action, true
		},
	})
}
