package builtin

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

var IOFuncers = []shapes.Funcer{
	shapes.Funcer{
		Summary:           "Groups lines into blocks separated by empty lines.",
		InputType:         types.NewArr(types.Str{}),
		InputDescription:  "an array of consecutive lines",
		Name:              "blocks",
		Params:            nil,
		OutputType:        types.NewArr(types.NewArr(types.Str{})),
		OutputDescription: "an array of arrays of lines, each representing a block",
		Notes:             "Each empty line in the input marks the end of a block. Blocks can be empty. The empty lines themselves are not included.",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		},
		IDs: nil,
		Examples: []shapes.Example{
			{`["a", "b", "", "c", "d", "e", "f", ""] blocks`, `Arr<Arr<Str...>...>`, `[["a", "b"], ["c", "d", "e", "f"]]`, nil},
			{`["a", ""] blocks`, `Arr<Arr<Str...>...>`, `[["a"]]`, nil},
			{`["a"] blocks`, `Arr<Arr<Str...>...>`, `[["a"]]`, nil},
			{`["", "a"] blocks`, `Arr<Arr<Str...>...>`, `[[], ["a"]]`, nil},
			{`["a", "", "", "", "b"] blocks`, `Arr<Arr<Str...>...>`, `[["a"], [], [], ["b"]]`, nil},
		},
	},
	shapes.SimpleFuncer(
		"Writes to STDERR.",
		types.NewVar("A", types.Any{}),
		"any value",
		"err",
		nil,
		types.NewVar("A", types.Any{}),
		"the same value",
		"Identity function with the side effect of writing a string representation of the value to STDERR, followed by a line break.",
		func(inputValue states.Value, args []states.Value) (states.Value, error) {
			str, err := inputValue.Str()
			if err != nil {
				return nil, err
			}
			fmt.Fprintln(os.Stderr, str)
			return inputValue, nil
		},
		nil,
	),
	shapes.SimpleFuncer(
		"Writes to STDERR with a custom line end.",
		types.NewVar("A", types.Any{}),
		"any value",
		"err",
		[]*params.Param{
			params.SimpleParam("end", "the line end to use", types.Str{}),
		},
		types.NewVar("A", types.Any{}),
		"the same value",
		"Identity function with the side effect of writing a string representation of the value to STDERR, followed by a the specified line end.",
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
		nil,
	),
	shapes.SimpleFuncer(
		"Reads from STDIN.",
		types.Any{},
		"any value (is ignored)",
		"in",
		nil,
		types.Reader{},
		"a Reader representing STDIN",
		"",
		func(inputValue states.Value, argValues []states.Value) (states.Value, error) {
			return states.ReaderValue{Reader: os.Stdin}, nil
		},
		nil,
	),
	shapes.Funcer{
		Summary:           "Reads JSON values from a stream",
		InputType:         types.Reader{},
		InputDescription:  "a Reader",
		Name:              "json",
		Params:            nil,
		OutputType:        types.AnyArr,
		OutputDescription: "array of data structures as they appear in the stream",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		},
		IDs:      nil,
		Examples: nil,
	},
	shapes.Funcer{
		Summary:           "Reads a stream line-by-line",
		InputType:         types.Reader{},
		InputDescription:  "a Reader",
		Name:              "lines",
		Params:            nil,
		OutputType:        types.NewArr(types.Str{}),
		OutputDescription: "an array of lines, without the line-break character",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		},
		IDs: nil,
		Examples: []shapes.Example{
			{`"abc\nde\n\nf" reader lines`, `Arr<Str...>`, `["abc", "de", "", "f"]`, nil},
		},
	},
	shapes.SimpleFuncer(
		"Writes to STDOUT.",
		types.NewVar("A", types.Any{}),
		"any value",
		"out",
		nil,
		types.NewVar("A", types.Any{}),
		"the same value",
		"Identity function with the side effect of writing a string representation of the value to STDERR, followed by a line break.",
		func(inputValue states.Value, args []states.Value) (states.Value, error) {
			str, err := inputValue.Str()
			if err != nil {
				return nil, err
			}
			fmt.Println(str)
			return inputValue, nil
		},
		nil,
	),
	shapes.SimpleFuncer(
		"Writes to STDOUT with a custom line end.",
		types.NewVar("A", types.Any{}),
		"any value",
		"out",
		[]*params.Param{
			params.SimpleParam("end", "", types.Str{}),
		},
		types.NewVar("A", types.Any{}),
		"the same value",
		"Identity function with the side effect of writing a string representation of the value to STDOUT, followed by a line break.",
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
		nil,
	),
	shapes.SimpleFuncer(
		"Creates a Reader from a Str.",
		types.Str{},
		"a string",
		"reader",
		nil,
		types.Reader{},
		"a Reader from which the input can be read",
		"",
		func(inputValue states.Value, args []states.Value) (states.Value, error) {
			return states.ReaderValue{strings.NewReader(string(inputValue.(states.StrValue)))}, nil
		},
		nil,
	),
}
