package builtin

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

var ObjFuncers = []shapes.Funcer{
	shapes.Funcer{
		Summary:          "Combines the items of two objects into one.",
		InputType:        types.AnyObj,
		InputDescription: "an object",
		Name:             "+",
		Params: []*params.Param{
			params.SimpleParam("other", "another object", types.AnyObj),
		},
		OutputType:        types.AnyObj,
		OutputDescription: "an object with the items of the input and other",
		Notes:             "Where the input and other share keys, the values of other are used in the output.",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			res := states.ObjValue{}
			for k, v := range inputState.Value.(states.ObjValue) {
				res[k] = v
			}
			arg, err := args[0](inputState.Clear(), nil).EvalObj()
			if err != nil {
				return states.ThunkFromError(err)
			}
			for k, v := range arg {
				res[k] = v
			}
			return states.ThunkFromValue(res)
		},
		IDs: nil,
		Examples: []shapes.Example{
			{`{a: 1} +{b: 2}`, `Obj<>`, `{a: 1, b: 2}`, nil},
			{`{a: 1, b: 2} +{b: 3}`, `Obj<>`, `{a: 1, b: 3}`, nil},
		},
	},
	shapes.Funcer{
		Summary: "Retrieves the value of a given property.",
		InputType: types.Obj{
			Props: map[string]types.Type{},
			Rest:  types.NewVar("A", types.Any{}),
		},
		InputDescription: "an object",
		Name:             "get",
		Params: []*params.Param{
			params.SimpleParam("prop", "a property key", types.NewUnion(types.Str{}, types.Num{})),
		},
		OutputType:        types.NewVar("A", types.Any{}),
		OutputDescription: "the value associated in the input object with the given property key",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			inputValue := inputState.Value.(states.ObjValue)
			val, err := args[0](inputState.Clear(), nil).Eval()
			if err != nil {
				return states.ThunkFromError(err)
			}
			prop, err := val.Str()
			if err != nil {
				return states.ThunkFromError(err)
			}
			thunk, ok := inputValue[prop]
			if !ok {
				return states.ThunkFromError(errors.ValueError(
					errors.Code(errors.NoSuchProperty),
					errors.GotValue(states.StrValue(prop)),
					errors.Pos(pos),
				))
			}
			val, err = thunk.Eval()
			if err != nil {
				return states.ThunkFromError(err)
			}
			return states.ThunkFromValue(val)
		},
		IDs: nil,
		Examples: []shapes.Example{
			{`{} get("a")`, `Void`, ``, errors.TypeError(
				errors.Code(errors.VoidProgram),
			)},
			{`{a: 1} get("a")`, `Num`, `1`, nil},
			{`{a: 1, b: "hey"} get("a")`, `Num|Str`, `1`, nil},
			{`{a: 1, b: "hey", c: false} get("a")`, `Num|Str|Bool`, `1`, nil},
			{`{1: "a"} get(1)`, `Str`, `"a"`, nil},
			{`{1.5: "a"} get(1.5)`, `Str`, `"a"`, nil},
			{`{b: 1} get("a")`, `Num`, ``, errors.ValueError(
				errors.Code(errors.NoSuchProperty),
				errors.GotValue(states.StrValue("a")),
			)},
		},
	},
	shapes.Funcer{
		Summary: "Checks for the presence of a given property.",
		InputType: types.Obj{
			Props: map[string]types.Type{},
			Rest:  types.NewVar("A", types.Any{}),
		},
		InputDescription: "an object",
		Name:             "has",
		Params: []*params.Param{
			params.SimpleParam("prop", "a property key", types.NewUnion(types.Str{}, types.Num{})),
		},
		OutputType:        types.Bool{},
		OutputDescription: "true if the object has the given property, false if not",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			inputValue := inputState.Value.(states.ObjValue)
			val, err := args[0](inputState.Clear(), nil).Eval()
			if err != nil {
				return states.ThunkFromError(err)
			}
			prop, err := val.Str()
			if err != nil {
				return states.ThunkFromError(err)
			}
			_, ok := inputValue[prop]
			return states.ThunkFromValue(states.BoolValue(ok))
		},
		IDs: nil,
		Examples: []shapes.Example{
			{`{} has("a")`, `Bool`, `false`, nil},
			{`{a: 1} has("a")`, `Bool`, `true`, nil},
		},
	},
	shapes.Funcer{
		Summary: "Iterates over the properties together with the values.",
		InputType: types.Obj{
			Props: map[string]types.Type{},
			Rest:  types.NewVar("A", types.Any{}),
		},
		InputDescription: "an object",
		Name:             "items",
		Params:           nil,
		OutputType: types.NewArr(
			types.NewTup([]types.Type{
				types.Str{},
				types.NewVar("A", types.Any{}),
			}),
		),
		OutputDescription: "an array of tuples of the properties and associated values of the input object",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			inputValue := inputState.Value.(states.ObjValue)
			c := make(chan *states.Thunk)
			go func() {
				for prop, thk := range inputValue {
					val, err := thk.Eval()
					if err != nil {
						c <- states.ThunkFromError(err)
						return
					}
					item := states.NewArrValue([]states.Value{
						states.StrValue(prop),
						val,
					})
					c <- states.ThunkFromValue(item)
				}
				c <- states.ThunkFromValue(nil)
			}()
			return states.ThunkFromChannel(c)
		},
		IDs: nil,
		Examples: []shapes.Example{
			{`{a: 1, b: 2} items sortBy(@0, <)`, `Arr<Arr<Str, Num>...>`, `[["a", 1], ["b", 2]]`, nil},
			{`{} items`, `Arr<Arr<Str, Void>...>`, `[]`, nil},
		},
	},
	shapes.Funcer{
		Summary: "Lists the properties of an object.",
		InputType: types.Obj{
			Props: map[string]types.Type{},
			Rest:  types.NewVar("A", types.Any{}),
		},
		InputDescription:  "an object",
		Name:              "props",
		Params:            nil,
		OutputType:        types.NewArr(types.Str{}),
		OutputDescription: "all the property keys of the object, as an array",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			inputValue := inputState.Value.(states.ObjValue)
			c := make(chan *states.Thunk)
			go func() {
				for prop := range inputValue {
					c <- states.ThunkFromValue(
						states.StrValue(prop),
					)
				}
				c <- states.ThunkFromValue(nil)
			}()
			return states.ThunkFromChannel(c)
		},
		IDs: nil,
		Examples: []shapes.Example{
			{`{a: 1, b: 2} props sort`, `Arr<Str...>`, `["a", "b"]`, nil},
			{`{} props`, `Arr<Str...>`, `[]`, nil},
		},
	},
	shapes.Funcer{
		Summary: "Lists the values of an object.",
		InputType: types.Obj{
			Props: map[string]types.Type{},
			Rest:  types.NewVar("A", types.Any{}),
		},
		InputDescription:  "an object",
		Name:              "values",
		Params:            nil,
		OutputType:        types.NewArr(types.NewVar("A", types.Any{})),
		OutputDescription: "all the property values of the object, as an array",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			inputValue := inputState.Value.(states.ObjValue)
			c := make(chan *states.Thunk)
			go func() {
				for _, thk := range inputValue {
					c <- thk
				}
				c <- states.ThunkFromValue(nil)
			}()
			return states.ThunkFromChannel(c)
		}, IDs: nil,
		Examples: []shapes.Example{
			{`{a: 1, b: 2} values sort`, `Arr<Num...>`, `[1, 2]`, nil},
			{`{} values`, `Arr<>`, `[]`, nil},
		},
	},
	shapes.Funcer{
		Summary: "Removes a property from an object.",
		InputType: types.Obj{
			Props: map[string]types.Type{},
			Rest:  types.NewVar("A", types.Any{}),
		},
		InputDescription: "an object",
		Name:             "without",
		Params: []*params.Param{
			params.SimpleParam("prop", "a property key", types.Str{}),
		},
		OutputType: types.Obj{
			Props: map[string]types.Type{},
			Rest:  types.NewVar("A", types.Any{}),
		},
		OutputDescription: "the input object, but with the specified property removed",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			inputValue := inputState.Value.(states.ObjValue)
			prop, err := args[0](inputState.Clear(), nil).EvalStr()
			if err != nil {
				return states.ThunkFromError(err)
			}
			res := states.ObjValue{}
			for k, v := range inputValue {
				if k != prop {
					res[k] = v
				}
			}
			return states.ThunkFromValue(res)
		},
		IDs: nil,
		Examples: []shapes.Example{
			{`{a: 1, b: 2} without("b")`, `Obj<Num>`, `{a: 1}`, nil},
			{`{a: 1, b: 2} without("c")`, `Obj<Num>`, `{a: 1, b: 2}`, nil},
		},
	},
}
