package patterns

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type TypePattern struct {
	Pos  lexer.Position
	Type types.Type
}

func (p TypePattern) Typecheck(inputShape shapes.Shape) (shapes.Shape, types.Type, Matcher, error) {
	if !inputShape.Type.Subsumes(p.Type) {
		return shapes.Shape{}, nil, nil, errors.E(
			errors.Code(errors.ImpossibleMatch),
			errors.Pos(p.Pos),
			errors.WantType(inputShape.Type),
			errors.GotType(p.Type),
		)
	}
	outputShape := shapes.Shape{
		Type:  p.Type,
		Stack: inputShape.Stack,
	}
	restType := types.Complement(inputShape.Type, p.Type)
	var check func(values.Value) bool
	switch t := p.Type.(type) {
	case types.NullType:
		check = func(v values.Value) bool {
			_, ok := v.(values.NullValue)
			return ok
		}
	case types.ReaderType:
		check = func(v values.Value) bool {
			_, ok := v.(values.ReaderValue)
			return ok
		}
	case types.BoolType:
		check = func(v values.Value) bool {
			_, ok := v.(values.BoolValue)
			return ok
		}
	case types.NumType:
		check = func(v values.Value) bool {
			_, ok := v.(values.NumValue)
			return ok
		}
	case types.StrType:
		check = func(v values.Value) bool {
			_, ok := v.(values.StrValue)
			return ok
		}
	case *types.SeqType:
		check = func(v values.Value) bool {
			u, ok := v.(values.SeqValue)
			return ok && t.ElType.Subsumes(u.ElementType)
		}
	default:
		panic("invalid type pattern")
	}
	matcher := func(inputState states.State) (*states.VariableStack, bool) {
		if check(inputState.Value) {
			return inputState.Stack, true
		}
		return nil, false
	}
	return outputShape, restType, matcher, nil
}
