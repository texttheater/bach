package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
)

type CallExpression struct {
	Pos  lexer.Position
	Name string
	Args []Expression
}

func (x CallExpression) Typecheck(inputShape shapes.Shape, params []*shapes.Parameter) (shapes.Shape, states.Action, error) {
	// Go down the function stack and find the function invoked by this
	// call
	stack := inputShape.Stack
	for {
		// Reached bottom of stack without finding a matching function
		if stack == nil {
			return shapes.Shape{}, nil, errors.E(
				errors.Code(errors.NoSuchFunction),
				errors.Pos(x.Pos),
				errors.InputType(inputShape.Type),
				errors.Name(x.Name),
				errors.NumParams(len(x.Args)+len(params)),
			)
		}
		// Try the funcer on top of the stack
		funcer := stack.Head
		funParams, funOutputType, funAction, ok := funcer(inputShape.Type, x.Name, len(x.Args)+len(params))
		if !ok {
			stack = stack.Tail
			continue
		}
		// Prepare action:
		action := funAction
		// Check function params filled by this call
		for i := 0; i < len(x.Args); i++ {
			param := funParams[i]
			argInputShape := shapes.Shape{param.InputType, inputShape.Stack}
			argOutputShape, argAction, err := x.Args[i].Typecheck(argInputShape, param.Params)
			if err != nil {
				return shapes.Shape{}, nil, err
			}
			if !param.OutputType.Subsumes(argOutputShape.Type) {
				return shapes.Shape{}, nil, errors.E(
					errors.Code(errors.ArgHasWrongOutputType),
					errors.Pos(x.Pos),
					errors.ArgNum(i),
					errors.WantType(param.OutputType),
					errors.GotType(argOutputShape.Type),
				)
			}
			action = action.SetArg(argAction)
		}
		// Check function params *not* filled by this call (thus left
		// to function to call)
		for i := 0; i < len(params); i++ {
			if !params[i].Subsumes(*funParams[len(x.Args)+i]) {
				return shapes.Shape{}, nil, errors.E(
					errors.Code(errors.ParamDoesNotMatch),
					errors.Pos(x.Pos),
					errors.ParamNum(i),
					errors.WantParam(params[i]),
					errors.GotParam(funParams[len(x.Args)+i]),
				)
			}
		}
		// Return result
		return shapes.Shape{funOutputType, inputShape.Stack}, action, nil
	}
}
