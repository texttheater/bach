package shapes

import (
	"fmt"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type Funcer struct {
	InputType  types.Type
	Name       string
	Params     []*params.Param // TODO does this need to be a slice of pointers?
	OutputType types.Type
	Kernel     RegularKernel
	IDs        *states.IDStack
}

type SimpleKernel func(inputValue states.Value, argValues []states.Value) (states.Value, error)

func SimpleFuncer(wantInputType types.Type, wantName string, pars []*params.Param, outputType types.Type, simpleKernel SimpleKernel) Funcer {
	// make regular kernel from simple kernel
	regularKernel := func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		argValues := make([]states.Value, len(pars))
		for i, arg := range args {
			val, err := arg(inputState.Clear(), nil).Eval()
			if err != nil {
				return states.ThunkFromError(err)
			}
			argValues[i] = val
		}
		value, err := simpleKernel(inputState.Value, argValues)
		if err != nil {
			return states.ThunkFromError(err)

		}
		return states.ThunkFromState(states.State{
			Value:     value,
			Stack:     inputState.Stack,
			TypeStack: inputState.TypeStack,
		})

	}
	// return
	return RegularFuncer(wantInputType, wantName, pars, outputType, regularKernel, nil)
}

func VariableFuncer(id any, name string, varType types.Type) Funcer {
	kernel := func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		stack := inputState.Stack
		for stack != nil {
			if stack.Head.ID == id {
				val, err := stack.Head.Action(states.InitialState, nil).Eval()
				if err != nil {
					return states.ThunkFromError(err)
				}
				return states.ThunkFromState(states.State{
					Value:     val,
					Stack:     inputState.Stack,
					TypeStack: inputState.TypeStack,
				})

			}
			stack = stack.Tail
		}
		panic(fmt.Sprintf("variable %s not found", name))
	}
	return RegularFuncer(types.Any{}, name, nil, varType, kernel, &states.IDStack{
		Head: id,
	})
}

type RegularKernel func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk

func RegularFuncer(wantInputType types.Type, wantName string, pars []*params.Param, outputType types.Type, kernel RegularKernel, ids *states.IDStack) Funcer {
	return Funcer{
		InputType:  wantInputType,
		Name:       wantName,
		Params:     pars,
		OutputType: outputType,
		Kernel:     kernel,
		IDs:        ids,
	}
}
