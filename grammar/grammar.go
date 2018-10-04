package grammar

import (
	//"fmt"
	//"os"
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
		`|(?P<Op2Number>(?:==|<=|>=)(?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+))` +
		`|(?P<NameLpar>(?:[+\-*/%<>=]|==|<=|>=)\()` +
		`|(?P<Comma>,)` +
		`|(?P<Rpar>\))`,
))

///////////////////////////////////////////////////////////////////////////////

type Composition struct {
	Components []*Component `{ @@ }`
}

func (c *Composition) ast() ast.Expression {
	var e ast.Expression
	e = &ast.IdentityExpression{}
	for _, comp := range c.Components {
		e = &ast.CompositionExpression{e, comp.ast()}
	}
	return e
}

///////////////////////////////////////////////////////////////////////////////

type Component struct {
	Pos     lexer.Position
	Number  *float64 `  @Number`
	NFFCall *NFFCall `| @@`
}

func (c *Component) ast() ast.Expression {
	if c.Number != nil {
		return &ast.NumberExpression{c.Pos, *c.Number}
	}
	if c.NFFCall != nil {
		return c.NFFCall.ast()
	}
	panic("invalid component")
}

///////////////////////////////////////////////////////////////////////////////

type NFFCall struct {
	Op1Number   *Op1Number   `  @Op1Number`
	Op2Number   *Op2Number   `| @Op2Number`
	NameArglist *NameArglist `| @@`
}

func (c *NFFCall) ast() ast.Expression {
	if c.Op1Number != nil {
		return c.Op1Number.ast()
	}
	if c.Op2Number != nil {
		return c.Op2Number.ast()
	}
	if c.NameArglist != nil {
		return c.NameArglist.ast()
	}
	panic("invalid NFF call")
}

///////////////////////////////////////////////////////////////////////////////

type Op1Number struct {
	Pos    lexer.Position
	Op     string
	Number float64
}

func (c *Op1Number) Capture(values []string) error {
	c.Op = string(values[0][:1])
	f, err := strconv.ParseFloat(values[0][1:], 64)
	if err != nil {
		return err
	}
	c.Number = f
	return nil
}

func (c *Op1Number) ast() ast.Expression {
	return &ast.NFFCallExpression{c.Pos, c.Op, []ast.Expression{&ast.NumberExpression{c.Pos, c.Number}}}
}

///////////////////////////////////////////////////////////////////////////////

type Op2Number struct {
	Pos    lexer.Position
	Op     string
	Number float64
}

func (c *Op2Number) Capture(values []string) error {
	c.Op = string(values[0][:2])
	f, err := strconv.ParseFloat(values[0][2:], 64)
	if err != nil {
		return err
	}
	c.Number = f
	return nil
}

func (c *Op2Number) ast() ast.Expression {
	return &ast.NFFCallExpression{c.Pos, c.Op, []ast.Expression{&ast.NumberExpression{c.Pos, c.Number}}}
}

///////////////////////////////////////////////////////////////////////////////

type NameArglist struct {
	Pos      lexer.Position
	NameLpar *NameLpar      `@NameLpar`
	Arg      *Composition   `@@`
	Args     []*Composition `{ "," @@ } ")"`
}

func (c *NameArglist) ast() ast.Expression {
	args := make([]ast.Expression, len(c.Args)+1)
	args[0] = c.Arg.ast()
	for i, Arg := range c.Args {
		args[i+1] = Arg.ast()
	}
	return &ast.NFFCallExpression{c.Pos, c.NameLpar.Name, args}
}

///////////////////////////////////////////////////////////////////////////////

type NameLpar struct {
	Pos  lexer.Position
	Name string
}

func (c *NameLpar) Capture(values []string) error {
	c.Name = values[0][:len(values[0])-1]
	return nil
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
