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
		expressions.RegularFuncer(
			&types.ArrType{types.StrType{}},
			"blocks",
			nil,
			&types.ArrType{&types.ArrType{types.StrType{}}},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				var nextBlock func(lines *states.ArrValue) (*states.ArrValue, *states.ArrValue, error)
				nextBlock = func(lines *states.ArrValue) (*states.ArrValue, *states.ArrValue, error) {
					if lines == nil {
						return nil, nil, nil
					}
					head := lines.Head.(states.StrValue)
					res := lines.Tail.Eval()
					if res.Error != nil {
						return nil, nil, res.Error
					}
					tail := res.Value.(*states.ArrValue)
					if head == "" {
						return nil, tail, nil
					}
					next, rest, err := nextBlock(tail)
					if err != nil {
						return nil, nil, err
					}
					return &states.ArrValue{
						Head: head,
						Tail: states.ThunkFromValue(next),
					}, rest, nil
				}
				lines := inputState.Value.(*states.ArrValue)
				var iter func() (states.Value, bool, error)
				iter = func() (states.Value, bool, error) {
					if lines == nil {
						return nil, false, nil
					}
					var next *states.ArrValue
					var err error
					next, lines, err = nextBlock(lines)
					if err != nil {
						return nil, false, err
					}
					return next, true, nil
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
				str, err := inputValue.Str()
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
				str, err := inputValue.Str()
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
				str, err := inputValue.Str()
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
				str, err := inputValue.Str()
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
