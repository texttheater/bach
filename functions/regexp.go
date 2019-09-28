package functions

import (
	"regexp"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type RegexpMatchExpression struct {
	Pos    lexer.Position
	Regexp *regexp.Regexp
}

func (x RegexpMatchExpression) Position() lexer.Position {
	return x.Pos
}

func (x RegexpMatchExpression) Typecheck(inputShape Shape, params []*Parameter) (Shape, states.Action, error) {
	if !(types.StrType{}).Subsumes(inputShape.Type) {
		return Shape{}, nil, errors.E(
			errors.Code(errors.RegexpWantsString),
			errors.Pos(x.Pos),
			errors.WantType(types.StrType{}),
			errors.GotType(inputShape.Type),
		)
	}
	panic("not implemented yet")
}
