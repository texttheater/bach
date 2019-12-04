package builtin

import (
	"bufio"
	"fmt"
	"os"

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
			func(inputState states.State, args []states.Action) *states.Thunk {
				reader := inputState.Value.(states.ReaderValue)
				scanner := bufio.NewScanner(reader.Reader)
				var next func() (states.Value, bool, error)
				next = func() (states.Value, bool, error) {
					ok := scanner.Scan()
					if !ok {
						return nil, false, nil
					}
					return states.StrValue(scanner.Text()), true, nil
				}
				return states.ThunkFromIter(next)
			},
		),
		functions.SimpleFuncer(
			types.TypeVariable{"$"},
			"out",
			nil,
			types.TypeVariable{"$"},
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
			types.TypeVariable{"$"},
			"out",
			[]types.Type{
				types.StrType{},
			},
			types.TypeVariable{"$"},
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
			types.TypeVariable{"$"},
			"err",
			nil,
			types.TypeVariable{"$"},
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
			types.TypeVariable{"$"},
			"err",
			[]types.Type{
				types.StrType{},
			},
			types.TypeVariable{"$"},
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
