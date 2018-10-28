package grammar

import (
	// "fmt"
	// "os"
	"strconv"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/ast"
	"github.com/texttheater/bach/errors"
)

///////////////////////////////////////////////////////////////////////////////

var LexerDefinition = lexer.Must(lexer.Regexp(
	`([\s]+)` +
		`|(?P<Number>(?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+))` +
		`|(?P<String>"(?:\\.|[^"])*")` +
		`|(?P<Op1Number>[+\-*/%<>](?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+))` +
		`|(?P<Op2Number>(?:==|<=|>=)(?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+))` +
		`|(?P<Op1Name>[+\-*/%<>](?:[+\-*/%<>]|==|<=|>=|[\p{L}_][\p{L}_0-9]*))` +
		`|(?P<Op2Name>(?:==|<=|>=)(?:[+\-*/%<>=]|==|<=|>=|[\p{L}_][\p{L}_0-9]*))` +
		`|(?P<Assignment>=(?:[+\-*/%<>=]|==|<=|>=|[\p{L}_][\p{L}_0-9]*))` +
		`|(?P<NameLpar>(?:[+\-*/%<>=]|==|<=|>=|[\p{L}_][\p{L}_0-9]*)\()` +
		`|(?P<Name>(?:[+\-*/%<>=]|==|<=|>=|[\p{L}_][\p{L}_0-9]*))` +
		`|(?P<Comma>,)` +
		`|(?P<Rpar>\))` +
		`|(?P<Keyword>for|def|as|ok)`, // these will be scanned as Names, but mapped to Keyword tokens by name2keyword (see below)
))

///////////////////////////////////////////////////////////////////////////////

type Composition struct {
	Pos        lexer.Position
	Component  *Component   `@@`
	Components []*Component `{ @@ }`
}

func (g *Composition) ast() ast.Expression {
	pos := g.Component.Pos
	e := g.Component.ast()
	for _, comp := range g.Components {
		e = &ast.CompositionExpression{pos, e, comp.ast()}
	}
	return e
}

///////////////////////////////////////////////////////////////////////////////

type Component struct {
	Pos        lexer.Position
	Number     *float64    `  @Number`
	String     *string     `| @String`
	Call       *Call       `| @@`
	Assignment *Assignment `| @Assignment`
	Definition *Definition `| @@`
}

func (g *Component) ast() ast.Expression {
	if g.Number != nil {
		return &ast.NumberExpression{g.Pos, *g.Number}
	}
	if g.String != nil {
		return &ast.StringExpression{g.Pos, *g.String}
	}
	if g.Call != nil {
		return g.Call.ast()
	}
	if g.Assignment != nil {
		return g.Assignment.ast()
	}
	panic("invalid component")
}

///////////////////////////////////////////////////////////////////////////////

type Call struct {
	Pos         lexer.Position
	Op1Number   *Op1Number   `  @Op1Number`
	Op2Number   *Op2Number   `| @Op2Number`
	Op1Name     *Op1Name     `| @Op1Name`
	Op2Name     *Op2Name     `| @Op2Name`
	NameArglist *NameArglist `| @@`
	Name        *string      `| @Name`
}

func (g *Call) ast() ast.Expression {
	if g.Op1Number != nil {
		return g.Op1Number.ast()
	}
	if g.Op2Number != nil {
		return g.Op2Number.ast()
	}
	if g.Op1Name != nil {
		return g.Op1Name.ast()
	}
	if g.Op2Name != nil {
		return g.Op2Name.ast()
	}
	if g.NameArglist != nil {
		return g.NameArglist.ast()
	}
	if g.Name != nil {
		return &ast.CallExpression{g.Pos, *g.Name, []ast.Expression{}}
	}
	panic("invalid call")
}

///////////////////////////////////////////////////////////////////////////////

type Op1Number struct {
	Pos    lexer.Position
	Op     string
	Number float64
}

func (g *Op1Number) Capture(values []string) error {
	g.Op = string(values[0][:1])
	f, err := strconv.ParseFloat(values[0][1:], 64)
	if err != nil {
		return err
	}
	g.Number = f
	return nil
}

func (g *Op1Number) ast() ast.Expression {
	return &ast.CallExpression{g.Pos, g.Op, []ast.Expression{&ast.NumberExpression{g.Pos, g.Number}}}
}

///////////////////////////////////////////////////////////////////////////////

type Op2Number struct {
	Pos    lexer.Position
	Op     string
	Number float64
}

func (g *Op2Number) Capture(values []string) error {
	g.Op = string(values[0][:2])
	f, err := strconv.ParseFloat(values[0][2:], 64)
	if err != nil {
		return err
	}
	g.Number = f
	return nil
}

func (g *Op2Number) ast() ast.Expression {
	return &ast.CallExpression{g.Pos, g.Op, []ast.Expression{&ast.NumberExpression{g.Pos, g.Number}}}
}

///////////////////////////////////////////////////////////////////////////////

type Op1Name struct {
	Pos  lexer.Position
	Op   string
	Name string
}

func (g *Op1Name) Capture(values []string) error {
	g.Op = string(values[0][:1])
	g.Name = values[0][1:]
	return nil
}

func (g *Op1Name) ast() ast.Expression {
	return &ast.CallExpression{g.Pos, g.Op, []ast.Expression{&ast.CallExpression{g.Pos, g.Name, []ast.Expression{}}}}
}

///////////////////////////////////////////////////////////////////////////////

type Op2Name struct {
	Pos  lexer.Position
	Op   string
	Name string
}

func (g *Op2Name) Capture(values []string) error {
	g.Op = string(values[0][:2])
	g.Name = values[0][2:]
	return nil
}

func (g *Op2Name) ast() ast.Expression {
	return &ast.CallExpression{g.Pos, g.Op, []ast.Expression{&ast.CallExpression{g.Pos, g.Name, []ast.Expression{}}}}
}

///////////////////////////////////////////////////////////////////////////////

type Assignment struct {
	Pos  lexer.Position
	Name string
}

func (g *Assignment) Capture(values []string) error {
	g.Name = values[0][1:]
	return nil
}

func (g *Assignment) ast() ast.Expression {
	return &ast.AssignmentExpression{
		Pos: g.Pos,
		Name: g.Name,
	}
}

///////////////////////////////////////////////////////////////////////////////

type Definition struct {
	Pos lexer.Position ` "for" "def" "as" "ok" ` // TODO dummy syntax
}

///////////////////////////////////////////////////////////////////////////////

type NameArglist struct {
	Pos      lexer.Position
	NameLpar *NameLpar      `@NameLpar`
	Arg      *Composition   `@@`
	Args     []*Composition `{ "," @@ } ")"`
}

func (g *NameArglist) ast() ast.Expression {
	args := make([]ast.Expression, len(g.Args)+1)
	args[0] = g.Arg.ast()
	for i, Arg := range g.Args {
		args[i+1] = Arg.ast()
	}
	return &ast.CallExpression{g.Pos, g.NameLpar.Name, args}
}

///////////////////////////////////////////////////////////////////////////////

type NameLpar struct {
	Pos  lexer.Position
	Name string
}

func (g *NameLpar) Capture(values []string) error {
	g.Name = values[0][:len(values[0])-1]
	return nil
}

///////////////////////////////////////////////////////////////////////////////

func Parse(input string) (ast.Expression, error) {
	parser, err := participle.Build(&Composition{}, participle.Lexer(LexerDefinition), participle.Unquote(LexerDefinition, "String"), participle.Map(name2keyword))
	if err != nil {
		if lexerError, ok := err.(*lexer.Error); ok {
			return nil, errors.E("syntax", lexerError.Pos, lexerError.Message)
		}
		return nil, err
	}
	composition := &Composition{}
	err = parser.ParseString(input, composition)
	if err != nil {
		if lexerError, ok := err.(*lexer.Error); ok {
			return nil, errors.E("syntax", lexerError.Pos, lexerError.Message)
		}
		return nil, err
	}
	return composition.ast(), nil
}

func name2keyword(t lexer.Token) lexer.Token {
	if t.Type != LexerDefinition.Symbols()["Name"] {
		return t
	}
	if isKeyword(t.Value) {
		t.Type = LexerDefinition.Symbols()["Keyword"]
	}
	return t
}

func isKeyword(name string) bool {
	return name == "for" || name == "def" || name == "as" || name == "ok"
}

///////////////////////////////////////////////////////////////////////////////
