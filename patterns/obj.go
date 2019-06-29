package patterns

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/types"
)

type ObjPattern struct {
	Pos            lexer.Position
	PropPatternMap map[string]Pattern
}

func (f *ObjPattern) Typecheck(inputShape shapes.Shape) (shapes.Shape, types.Type, Matcher, error) {
	panic("not implemented")
}
