// Package functions implements functions from states to states.
//
// A state consists of a value and a stack of available variables (named
// values).
//
// Shapes are to states as types are to values. A shape consists of a type and
// a stack of available functions.
//
// Interpreting a Bach program involves assigning each expression an input
// shape, a function and an output shape. The first expression in the program
// gets the initial shape, consisting of the Any type and a stack consisting
// only of builtin functions. The input shape of an expression and the
// expression together determine its function. The function and the input shape
// together determine its output shape. In a composition expression L R, the
// output shape of L is the input shape of R.
package functions

import (
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type Function struct {
	InputType    types.Type
	Name         string      // TODO namespaces
	FilledParams []*Function // TODO make stack instead?
	OpenParams   []*Parameter
	UpdateShape  func(inputShape Shape) Shape
	Evaluate      func(inputValue values.Value, args []*Function) values.Value
}

func (f *Function) SetArg(arg *Function) *Function {
	filledParams := make([]*Function, 0, len(f.FilledParams)+1)
	filledParams = append(filledParams, f.FilledParams...)
	filledParams = append(filledParams, arg)
	return &Function{
		f.InputType,
		f.Name,
		filledParams,
		f.OpenParams[1:],
		f.UpdateShape,
		f.Evaluate,
	}
}

func (f *Function) Apply(inputValue values.Value, outsideArgs []*Function) values.Value {
	args := make([]*Function, 0, len(f.FilledParams)+len(outsideArgs))
	args = append(args, f.FilledParams...)
	args = append(args, outsideArgs...)
	return f.Evaluate(inputValue, args)
}

func SimpleFunction(inputType types.Type, name string, argTypes []types.Type,
	outputType types.Type,
	eval func(values.Value, []values.Value) values.Value) *Function {
	params := make([]*Parameter, 0, len(argTypes))
	for _, argType := range argTypes {
		params = append(params, &Parameter{
			InputType: &types.AnyType{},
			Parameters: nil,
			OutputType: argType,
		})
	}
	return &Function{
		InputType: inputType,
		Name: name,
		FilledParams: nil,
		OpenParams: params,
		UpdateShape: func(inputShape Shape) Shape {
			return Shape{
				outputType,
				inputShape.Stack,
			}
		},
		Evaluate: func(inputValue values.Value, args []*Function) values.Value {
			argValues := make([]values.Value, 0, len(args))
			for _, arg := range args {
				argValues = append(argValues, arg.Apply(inputValue , nil))
			}
			return eval(inputValue, argValues)
		},
	}
}
