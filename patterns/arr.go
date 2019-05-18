package patterns

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/types"
)

type ArrPattern struct {
	Pos             lexer.Position
	ElementPatterns []Pattern
}

func (p *ArrPattern) Typecheck(inputShape shapes.Shape) (shapes.Shape, types.Type, Matcher, error) {
	panic("not implemented")
}
