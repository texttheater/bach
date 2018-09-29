package grammar

import (
	"strconv"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/ast"
)

///////////////////////////////////////////////////////////////////////////////

var Lexer = lexer.Must(lexer.Regexp(
	`([\s]+)` +
	`|(?P<Number>(?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+))` +
	`|(?P<Op1Number>[+\-*/%<>=](?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+))` +
	`|(?P<Op2Number>(?:==|<=|>=)(?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+))`,
))

///////////////////////////////////////////////////////////////////////////////

type Component struct {
	Pos lexer.Position
	Number *float64 `@Number`
	Op1Number *Op1Number `|@Op1Number`
	Op2Number *Op2Number `|@Op2Number`
}

func (c *Component) ast() ast.Expression {
	if c.Number != nil {
		return ast.NumberExpression{c.Pos, *c.Number}
	}
	panic("invalid component")
}

///////////////////////////////////////////////////////////////////////////////

type Op1Number struct {
	Pos lexer.Position
	Op string
	Number *float64
}

func (o Op1Number) Capture(values []string) error {
	o.Op = string(values[0][:1])
	f, err := strconv.ParseFloat(values[0][1:], 64)
	if err != nil {
		return err
	}
	o.Number = &f
	return nil
}

///////////////////////////////////////////////////////////////////////////////

type Op2Number struct {
	Pos lexer.Position
	Op string
	Number *float64
}

func (o Op2Number) Capture(values []string) error {
	o.Op = string(values[0][:2])
	f, err := strconv.ParseFloat(values[0][2:], 64)
	if err != nil {
		return err
	}
	o.Number = &f
	return nil
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
