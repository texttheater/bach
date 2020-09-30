package builtin

import (
	"bufio"
	"fmt"
	"os"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initIO() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]expressions.Funcer{
		expressions.SimpleFuncer(
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
		expressions.RegularFuncer(
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
		expressions.SimpleFuncer(
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
		expressions.SimpleFuncer(
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
		expressions.SimpleFuncer(
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
		expressions.SimpleFuncer(
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
