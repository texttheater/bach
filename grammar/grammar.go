package grammar

import (
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/ast"
)

///////////////////////////////////////////////////////////////////////////////

var Lexer = lexer.Must(lexer.Regexp(
	`([\s]+)` +
	`|(?P<Number>(?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+))`,
))

///////////////////////////////////////////////////////////////////////////////

type Component struct {
	Pos lexer.Position
	Number *float64 `@Number`
}

func (c *Component) ast() ast.Expression {
	if c.Number != nil {
		return ast.NumberExpression{c.Pos, *c.Number}
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
