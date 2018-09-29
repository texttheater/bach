package grammar

import (
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/ast"
)

///////////////////////////////////////////////////////////////////////////////

var Lexer = lexer.Must(lexer.Regexp(
	`([\s]+)` +
	`|(?P<Float>(?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?))` +
	`|(?P<Int>(?:[1-9]\d*|0[0-7]*|0[xX][0-9a-fA-F]+))`,
))

///////////////////////////////////////////////////////////////////////////////

type Component struct {
	Pos lexer.Position
	Float *float64 `@Float`
	Int *int64 `|@Int`
}

func (c *Component) ast() ast.Expression {
	if c.Float != nil {
		return ast.NumberExpression{c.Pos, *c.Float}
	}
	if c.Int != nil {
		return ast.NumberExpression{c.Pos, float64(*c.Int)}
	}
	panic("invalid component")
}

///////////////////////////////////////////////////////////////////////////////

type Composition struct {
	Components []*Component `{ @@ }`
}

func (cc *Composition) ast() ast.Expression {
	var e ast.Expression
	e = ast.IdentityExpression{}
	for _, c := range cc.Components {
		e = ast.CompositionExpression{e, c.ast()}
	}
	return e
}

///////////////////////////////////////////////////////////////////////////////

func Parse(input string) (ast.Expression, error) {
	parser, err := participle.Build(&Composition{}, participle.Lexer(Lexer))
	if err != nil {
		return nil, err
	}
	composition := &Composition{}
	err = parser.ParseString(input, composition)
	if err != nil {
		return nil, err
	}
	return composition.ast(), nil
}
