package builtin

import (
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initRegexp() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		func(gotInputShape functions.Shape, gotCall functions.CallExpression, gotParams []*functions.Parameter) (functions.Shape, states.Action, bool, error) {
			// match number of parameters
			if len(gotCall.Args)+len(gotParams) != 0 {
				return functions.Shape{}, nil, false, nil
			}
			// match name
			if gotCall.Name != "arr" {
				return functions.Shape{}, nil, false, nil
			}
			// match input type
			if !(types.StrType{}).Subsumes(gotInputShape.Type) {
				return functions.Shape{}, nil, false, nil
			}
			// create action
			action := states.Action(func(inputState states.State, args []states.Action) states.State {
				return args[0](inputState, nil)
			})
			// typecheck parameter
			var outputType types.Type
			if len(gotCall.Args) > 0 { // parameter set by the call
				argInputShape := functions.Shape{
					Type:  types.StrType{},
					Stack: gotInputShape.Stack,
				}
				argOutputShape, argAction, err := gotCall.Args[0].Typecheck(argInputShape, nil)
				if err != nil {
					return functions.Shape{}, nil, false, err
				}
				if !(types.ObjType{}).Subsumes(argOutputShape.Type) {
					return functions.Shape{}, nil, false, errors.E(
						errors.Code(errors.ArgHasWrongOutputType),
						errors.Pos(gotCall.Pos),
						errors.ArgNum(0),
						errors.WantType(types.ObjType{}),
						errors.GotType(argOutputShape.Type),
					)
				}
				action = action.SetArg(argAction)
				outputType = argOutputShape.Type
			} else { // parameter not set by the call
				wantParam := &functions.Parameter{
					InputType:  types.StrType{},
					Params:     nil,
					OutputType: types.ObjType{},
				}
				if !gotParams[0].Subsumes(wantParam) {
					return functions.Shape{}, nil, false, errors.E(
						errors.Code(errors.ParamDoesNotMatch),
						errors.Pos(gotCall.Pos),
						errors.ParamNum(1),
						errors.WantParam(gotParams[0]),
						errors.GotParam(wantParam),
					)
				}
				outputType = gotParams[0].OutputType
			}
			// create output shape
			outputShape := functions.Shape{
				Type:  outputType,
				Stack: gotInputShape.Stack,
			}
			return outputShape, action, true, nil
		},
	})
}
