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
				head, err := lines.Head.EvalStr()
				if err != nil {
					return nil, nil, err
				}
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
					Head: states.ThunkFromValue(states.StrValue(head)),
					Tail: states.ThunkFromValue(next),
				}, rest, nil
			}
			lines, err := inputState.Thunk.EvalArr()
			if err != nil {
				return states.ThunkFromError(err)
			}
			iter := func() (*states.Thunk, bool, error) {
				if lines == nil {
					return nil, false, nil
				}
				var next *states.ArrValue
				var err error
				next, lines, err = nextBlock(lines)
				if err != nil {
					return nil, false, err
				}
				return states.ThunkFromValue(next), true, nil
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
		types.Null{},
		"null",
		"Writes a string representation of the value to STDERR, followed by a line break.",
		func(inputThunk *states.Thunk, argThunks []*states.Thunk) *states.Thunk {
			str, err := inputThunk.EvalStr()
			if err != nil {
				return states.ThunkFromError(err)
			}
			fmt.Fprintln(os.Stderr, str)
			return inputThunk
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
		types.Null{},
		"null",
		"Writes a string representation of the value to STDERR, followed by the specified line end.",
		func(inputThunk *states.Thunk, argThunks []*states.Thunk) *states.Thunk {
			str, err := inputThunk.EvalStr()
			if err != nil {
				return states.ThunkFromError(err)
			}
			end, err := argThunks[0].EvalStr()
			fmt.Fprint(os.Stderr, str)
			fmt.Fprint(os.Stderr, end)
			return inputThunk
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
		func(inputThunk *states.Thunk, argThunks []*states.Thunk) *states.Thunk {
			return states.ThunkFromValue(states.ReaderValue{Reader: os.Stdin})
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
			reader, err := inputState.Thunk.EvalReader()
			if err != nil {
				return states.ThunkFromError(err)
			}
			dec := json.NewDecoder(reader)
			output := func() (*states.Thunk, bool, error) {
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
				return thunkFromData(o, pos), true, nil
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
		Kernel:            Lines,
		IDs:               nil,
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
		types.Null{},
		"null",
		"Writes a string representation of the value to STDERR, followed by a line break.",
		func(inputThunk *states.Thunk, argThunks []*states.Thunk) *states.Thunk {
			val, err := inputThunk.Eval()
			if err != nil {
				return states.ThunkFromError(err)
			}
			str, err := val.Str()
			if err != nil {
				return states.ThunkFromError(err)
			}
			fmt.Println(str)
			return states.ThunkFromValue(states.NullValue{})
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
		types.Null{},
		"null",
		"Writes a string representation of the value to STDOUT, followed by the specified line end.",
		func(inputThunk *states.Thunk, argThunks []*states.Thunk) *states.Thunk {
			str, err := inputThunk.EvalStr()
			if err != nil {
				return states.ThunkFromError(err)
			}
			end, err := argThunks[0].EvalStr()
			if err != nil {
				return states.ThunkFromError(err)
			}
			fmt.Print(str)
			fmt.Print(end)
			return states.ThunkFromValue(states.NullValue{})
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
		func(inputThunk *states.Thunk, argThunks []*states.Thunk) *states.Thunk {
			str, err := inputThunk.EvalStr()
			if err != nil {
				return states.ThunkFromError(err)
			}
			return states.ThunkFromValue(states.ReaderValue{
				Reader: strings.NewReader(str),
			})
		},
		nil,
	),
}

func Lines(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
	return states.ThunkFromFunc(func() *states.Thunk {
		reader, err := inputState.Thunk.EvalReader()
		if err != nil {
			return states.ThunkFromError(err)
		}
		scanner := bufio.NewScanner(reader)
		iter := func() (*states.Thunk, bool, error) {
			ok := scanner.Scan()
			if !ok {
				return nil, false, nil
			}
			return states.ThunkFromValue(states.StrValue(scanner.Text())), true, nil
		}
		return states.ThunkFromIter(iter)
	})
}
