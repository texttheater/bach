package builtin

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

var IOFuncers = []shapes.Funcer{

	shapes.Funcer{InputType: types.NewArr(types.Str{}), Name: "blocks", Params: nil, OutputType: types.NewArr(types.NewArr(types.Str{})), Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		iter := func() (states.Value, bool, error) {
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
	}, IDs: nil},

	shapes.SimpleFuncer("", types.NewVar("A", types.Any{}), "", "err", nil, types.NewVar("A", types.Any{}), "", "", func(inputValue states.Value, args []states.Value) (states.Value, error) {
		str, err := inputValue.Str()
		if err != nil {
			return nil, err
		}
		fmt.Fprintln(os.Stderr, str)
		return inputValue, nil
	}, nil),

	shapes.SimpleFuncer("", types.NewVar("A", types.Any{}), "", "err", []*params.Param{
		params.SimpleParam("message", "", types.Str{}),
	}, types.NewVar("A", types.Any{}), "", "", func(inputValue states.Value, args []states.Value) (states.Value, error) {
		str, err := inputValue.Str()
		if err != nil {
			return nil, err
		}
		end := string(args[0].(states.StrValue))
		fmt.Fprint(os.Stderr, str)
		fmt.Fprint(os.Stderr, end)
		return inputValue, nil
	}, nil),

	shapes.SimpleFuncer("", types.Any{}, "", "in", nil, types.Reader{}, "", "", func(inputValue states.Value, argValues []states.Value) (states.Value, error) {
		return states.ReaderValue{Reader: os.Stdin}, nil
	}, nil),

	shapes.Funcer{InputType: types.Reader{}, Name: "json", Params: nil, OutputType: types.Any{}, Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		reader := inputState.Value.(states.ReaderValue).Reader
		dec := json.NewDecoder(reader)
		output := func() (states.Value, bool, error) {
			if !dec.More() {
				return nil, false, nil
			}
			var o any
			err := dec.Decode(&o)
			if err != nil {
				return nil, false, errors.ValueError(
					errors.Pos(pos),
					errors.Code(errors.UnexpectedValue),
					errors.Message(err.Error()),
				)
			}
			val, err := thunkFromData(o, pos).Eval()
			if err != nil {
				return nil, false, err
			}
			return val, true, nil
		}
		return states.ThunkFromIter(output)
	}, IDs: nil},

	shapes.Funcer{InputType: types.Reader{}, Name: "lines", Params: nil, OutputType: types.NewArr(types.Str{}), Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		reader := inputState.Value.(states.ReaderValue)
		scanner := bufio.NewScanner(reader.Reader)
		iter := func() (states.Value, bool, error) {
			ok := scanner.Scan()
			if !ok {
				return nil, false, nil
			}
			return states.StrValue(scanner.Text()), true, nil
		}
		return states.ThunkFromIter(iter)
	}, IDs: nil},

	shapes.SimpleFuncer("", types.NewVar("A", types.Any{}), "", "out", nil, types.NewVar("A", types.Any{}), "", "", func(inputValue states.Value, args []states.Value) (states.Value, error) {
		str, err := inputValue.Str()
		if err != nil {
			return nil, err
		}
		fmt.Println(str)
		return inputValue, nil
	}, nil),

	shapes.SimpleFuncer("", types.NewVar("A", types.Any{}), "", "out", []*params.Param{
		params.SimpleParam("message", "", types.Str{}),
	}, types.NewVar("A", types.Any{}), "", "", func(inputValue states.Value, args []states.Value) (states.Value, error) {
		str, err := inputValue.Str()
		if err != nil {
			return nil, err
		}
		end := string(args[0].(states.StrValue))
		fmt.Print(str)
		fmt.Print(end)
		return inputValue, nil
	}, nil),
}
