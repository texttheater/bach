package builtin

import (
	"bufio"
	"fmt"
	"os"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initIO() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		functions.SimpleFuncer(
			types.AnyType{},
			"in",
			nil,
			types.ReaderType{},
			func(inputValue states.Value, argValues []states.Value) (states.Value, error) {
				return states.ReaderValue{
					os.Stdin,
				}, nil
			},
		),
		functions.RegularFuncer(
			types.ReaderType{},
			"lines",
			nil,
			&types.ArrType{types.StrType{}},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				reader := inputState.Value.(states.ReaderValue)
				scanner := bufio.NewScanner(reader.Reader)
				var iter func() (states.Value, bool, error)
				iter = func() (states.Value, bool, error) {
					ok := scanner.Scan()
					if !ok {
						return nil, false, nil
					}
					return states.StrValue(scanner.Text()), true, nil
				}
				return states.ThunkFromIter(iter)
			},
			nil,
		),
		functions.SimpleFuncer(
			types.TypeVariable{
				Name: "A",
			},
			"out",
			nil,
			types.TypeVariable{
				Name: "A",
			},
			func(inputValue states.Value, args []states.Value) (states.Value, error) {
				str, err := inputValue.Out()
				if err != nil {
					return nil, err
				}
				fmt.Println(str)
				return inputValue, nil
			},
		),
		functions.SimpleFuncer(
			types.TypeVariable{
				Name: "A",
			},
			"out",
			[]types.Type{
				types.StrType{},
			},
			types.TypeVariable{
				Name: "A",
			},
			func(inputValue states.Value, args []states.Value) (states.Value, error) {
				str, err := inputValue.Out()
				if err != nil {
					return nil, err
				}
				end := string(args[0].(states.StrValue))
				fmt.Print(str)
				fmt.Print(end)
				return inputValue, nil
			},
		),
		functions.SimpleFuncer(
			types.TypeVariable{
				Name: "A",
			},
			"err",
			nil,
			types.TypeVariable{
				Name: "A",
			},
			func(inputValue states.Value, args []states.Value) (states.Value, error) {
				str, err := inputValue.Out()
				if err != nil {
					return nil, err
				}
				fmt.Fprintln(os.Stderr, str)
				return inputValue, nil
			},
		),
		functions.SimpleFuncer(
			types.TypeVariable{
				Name: "A",
			},
			"err",
			[]types.Type{
				types.StrType{},
			},
			types.TypeVariable{
				Name: "A",
			},
			func(inputValue states.Value, args []states.Value) (states.Value, error) {
				str, err := inputValue.Out()
				if err != nil {
					return nil, err
				}
				end := string(args[0].(states.StrValue))
				fmt.Fprint(os.Stderr, str)
				fmt.Fprint(os.Stderr, end)
				return inputValue, nil
			},
		),
	})
}
