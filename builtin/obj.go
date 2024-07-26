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
		InputType: types.AnyObj,
		Name:      "+",
		Params: []*params.Param{
			params.SimpleParam("other", "", types.AnyObj),
		},
		OutputType: types.AnyObj,
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
	},
	shapes.Funcer{
		InputType: types.Obj{
			Props: map[string]types.Type{},
			Rest:  types.NewVar("A", types.Any{}),
		},
		Name: "get",
		Params: []*params.Param{
			params.SimpleParam("key", "", types.NewUnion(types.Str{}, types.Num{})),
		},
		OutputType: types.NewVar("A", types.Any{}),
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
	},
	shapes.Funcer{
		InputType: types.Obj{
			Props: map[string]types.Type{},
			Rest:  types.NewVar("A", types.Any{}),
		},
		Name: "has",
		Params: []*params.Param{
			params.SimpleParam("key", "", types.NewUnion(types.Str{}, types.Num{})),
		},
		OutputType: types.Bool{},
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
	},
	shapes.Funcer{
		InputType: types.Obj{
			Props: map[string]types.Type{},
			Rest:  types.NewVar("A", types.Any{}),
		},
		Name:   "items",
		Params: nil,
		OutputType: types.NewArr(
			types.NewTup([]types.Type{
				types.Str{},
				types.NewVar("A", types.Any{}),
			}),
		),
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
	},
	shapes.Funcer{
		InputType: types.Obj{
			Props: map[string]types.Type{},
			Rest:  types.NewVar("A", types.Any{}),
		},
		Name:       "props",
		Params:     nil,
		OutputType: types.NewArr(types.Str{}),
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
	},
	shapes.Funcer{
		InputType: types.Obj{
			Props: map[string]types.Type{},
			Rest:  types.NewVar("A", types.Any{}),
		},
		Name:       "values",
		Params:     nil,
		OutputType: types.NewArr(types.NewVar("A", types.Any{})),
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
	},
	shapes.Funcer{
		InputType: types.Obj{
			Props: map[string]types.Type{},
			Rest:  types.NewVar("A", types.Any{}),
		},
		Name: "without",
		Params: []*params.Param{
			params.SimpleParam("key", "", types.Str{}),
		},
		OutputType: types.Obj{
			Props: map[string]types.Type{},
			Rest:  types.NewVar("A", types.Any{}),
		},
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
	},
}
