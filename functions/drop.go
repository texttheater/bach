package functions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type DropExpression struct {
	Pos lexer.Position
}

func (x DropExpression) Position() lexer.Position {
	return x.Pos
}

func (x DropExpression) Typecheck(inputShape Shape, params []*Parameter) (Shape, states.Action, error) {
	outputShape := Shape{
		Type:  types.VoidType{},
		Stack: nil,
	}
	action := func(inputState states.State, args []states.Action) states.State {
		return states.State{
			Drop:  true,
			Value: nil,
			Stack: nil,
		}
	}
	return outputShape, action, nil
}
