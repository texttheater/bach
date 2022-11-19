package builtin

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initObj() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]expressions.Funcer{
		// for Obj<Any> +(Obj<Any>) Obj<Any>
		expressions.RegularFuncer(
			types.AnyObj,
			"+",
			[]*params.Param{
				params.SimpleParam(types.AnyObj),
			},
			types.AnyObj,
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
			nil,
		),
		// for Obj<<A>> get(Str|Str) <A>
		expressions.RegularFuncer(
			types.Obj{
				Props: map[string]types.Type{},
				Rest:  types.NewVar("A", types.Any{}),
			},
			"get",
			[]*params.Param{
				params.SimpleParam(types.NewUnion(types.Str{}, types.Num{})),
			},
			types.NewVar("A", types.Any{}),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				inputValue := inputState.Value.(states.ObjValue)
				val, err := args[0](inputState.Clear(), nil).Eval()
				if err != nil {
					return states.ThunkFromError(err)
				}
				prop, err := val.Str() // TODO ???
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
			nil,
		),
		// for Obj<<A>> has(Str|Num) Bool
		expressions.RegularFuncer(
			types.Obj{
				Props: map[string]types.Type{},
				Rest:  types.NewVar("A", types.Any{}),
			},
			"has",
			[]*params.Param{
				params.SimpleParam(types.NewUnion(types.Str{}, types.Num{})),
			},
			types.Bool{},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				inputValue := inputState.Value.(states.ObjValue)
				val, err := args[0](inputState.Clear(), nil).Eval()
				if err != nil {
					return states.ThunkFromError(err)
				}
				prop, err := val.Str() // TODO ???
				if err != nil {
					return states.ThunkFromError(err)
				}
				_, ok := inputValue[prop]
				return states.ThunkFromValue(states.BoolValue(ok))
			},
			nil,
		),
		// for Obj<<A>> props Arr<Str>
		expressions.RegularFuncer(
			types.Obj{
				Props: map[string]types.Type{},
				Rest:  types.NewVar("A", types.Any{}),
			},
			"props",
			nil,
			types.NewArr(types.Str{}),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
			nil,
		),
		// for Obj<<A>> values Arr<<A>>
		expressions.RegularFuncer(
			types.Obj{
				Props: map[string]types.Type{},
				Rest:  types.NewVar("A", types.Any{}),
			},
			"values",
			nil,
			types.NewArr(types.NewVar("A", types.Any{})),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				inputValue := inputState.Value.(states.ObjValue)
				c := make(chan *states.Thunk)
				go func() {
					for _, thk := range inputValue {
						c <- thk
					}
					c <- states.ThunkFromValue(nil)
				}()
				return states.ThunkFromChannel(c)
			},
			nil,
		),
	})
}
