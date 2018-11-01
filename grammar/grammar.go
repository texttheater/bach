package grammar

import (
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/ast"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
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
		`|(?P<Keyword>for|def|as|ok|Num)`, // these will be scanned as Names, but mapped to Keyword tokens by name2keyword (see below)
))

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
	return composition.Ast(), nil
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
	return name == "for" || name == "def" || name == "as" || name == "ok" || name == "Num"
}

///////////////////////////////////////////////////////////////////////////////

type Composition struct {
	Pos        lexer.Position
	Component  *Component   `@@`
	Components []*Component `{ @@ }`
}

func (g *Composition) Ast() ast.Expression {
	pos := g.Component.Pos
	e := g.Component.Ast()
	for _, comp := range g.Components {
		e = &ast.CompositionExpression{pos, e, comp.Ast()}
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

func (g *Component) Ast() ast.Expression {
	if g.Number != nil {
		return &ast.ConstantExpression{
			Pos:   g.Pos,
			Type:  &types.NumberType{},
			Value: &values.NumberValue{*g.Number},
		}
	}
	if g.String != nil {
		return &ast.ConstantExpression{
			Pos:   g.Pos,
			Type:  &types.StringType{},
			Value: &values.StringValue{*g.String},
		}
	}
	if g.Call != nil {
		return g.Call.Ast()
	}
	if g.Assignment != nil {
		return g.Assignment.Ast()
	}
	panic("invalid component")
}

///////////////////////////////////////////////////////////////////////////////
