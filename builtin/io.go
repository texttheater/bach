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
			types.Any{},
			"in",
			nil,
			types.Reader{},
			func(inputValue states.Value, argValues []states.Value) (states.Value, error) {
				return states.ReaderValue{Reader: os.Stdin}, nil
			},
		),
		expressions.RegularFuncer(
			types.Reader{},
			"lines",
			nil,
			types.NewArr(types.Str{}),
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
			types.NewArr(types.Str{}),
			"blocks",
			nil,
			types.NewArr(types.NewArr(types.Str{})),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				var nextBlock func(lines *states.ArrValue) (*states.ArrValue, *states.ArrValue, error)
				nextBlock = func(lines *states.ArrValue) (*states.ArrValue, *states.ArrValue, error) {
					if lines == nil {
						return nil, nil, nil
					}
					head := lines.Head.(states.StrValue)
					val, err := lines.Tail.Eval()
					if err != nil {
						return nil, nil, err
					}
					tail := val.(*states.ArrValue)
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
			types.NewVar("A", types.Any{}),
			"out",
			nil,
			types.NewVar("A", types.Any{}),
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
			types.NewVar("A", types.Any{}),
			"out",
			[]types.Type{
				types.Str{},
			},
			types.NewVar("A", types.Any{}),
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
			types.NewVar("A", types.Any{}),
			"err",
			nil,
			types.NewVar("A", types.Any{}),
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
			types.NewVar("A", types.Any{}),
			"err",
			[]types.Type{
				types.Str{},
			},
			types.NewVar("A", types.Any{}),
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
