package functions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type LeftExpression struct {
	Pos lexer.Position
}

func (x LeftExpression) Typecheck(inputShape Shape, params []*Parameter) (Shape, states.Action, error) {
	outputShape := Shape{
		Type: types.NewObjType(map[string]types.Type{
			"left": inputShape.Type,
		}),
		Stack: inputShape.Stack,
	}
	action := func(inputState states.State, args []states.Action) states.State {
		return states.State{
			Value: values.ObjValue(map[string]values.Value{
				"left": inputState.Value,
			}),
			Stack: inputState.Stack,
		}
	}
	return outputShape, action, nil
}

type RightExpression struct {
	Pos lexer.Position
}

func (x RightExpression) Typecheck(inputShape Shape, params []*Parameter) (Shape, states.Action, error) {
	outputShape := Shape{
		Type: types.NewObjType(map[string]types.Type{
			"right": inputShape.Type,
		}),
		Stack: inputShape.Stack,
	}
	action := func(inputState states.State, args []states.Action) states.State {
		return states.State{
			Value: values.ObjValue(map[string]values.Value{
				"right": inputState.Value,
			}),
			Stack: inputState.Stack,
		}
	}
	return outputShape, action, nil
}
