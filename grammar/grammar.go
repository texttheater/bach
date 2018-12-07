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
		`|(?P<Lbrack>\[)` +
		`|(?P<Rbrack>])` +
		// the following will be scanned as Name, but mapped to the appropriate token types by name2keyword (see below)
		`|(?P<Keyword>for|def|as|ok|if|then|elif|else)` +
		`|(?P<Type>Num|Str|Bool|Null|Any)`,
))

///////////////////////////////////////////////////////////////////////////////

func Parse(input string) (ast.Expression, error) {
	parser, err := participle.Build(&Composition{}, participle.Lexer(LexerDefinition), participle.Unquote(LexerDefinition, "String"), participle.Map(Name2keyword))
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

func Name2keyword(t lexer.Token) lexer.Token {
	if t.Type != LexerDefinition.Symbols()["Name"] {
		return t
	}
	if isTypeKeyword(t.Value) {
		t.Type = LexerDefinition.Symbols()["Type"]
	}
	if isKeyword(t.Value) {
		t.Type = LexerDefinition.Symbols()["Keyword"]
	}
	return t
}

func isTypeKeyword(name string) bool {
	return name == "Num" || name == "Str" || name == "Bool" ||
		name == "Null" || name == "Any"
}

func isKeyword(name string) bool {
	return name == "for" || name == "def" || name == "as" ||
		name == "ok" || name == "if" || name == "then" ||
		name == "elif" || name == "else"
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
	Pos         lexer.Position
	Number      *float64     `  @Number`
	String      *string      `| @String`
	Array       *Array       `| @@`
	Call        *Call        `| @@`
	Assignment  *Assignment  `| @Assignment`
	Definition  *Definition  `| @@`
	Conditional *Conditional `| @@`
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
	if g.Array != nil {
		return g.Array.Ast()
	}
	if g.Call != nil {
		return g.Call.Ast()
	}
	if g.Assignment != nil {
		return g.Assignment.Ast()
	}
	if g.Definition != nil {
		return g.Definition.Ast()
	}
	if g.Conditional != nil {
		return g.Conditional.Ast()
	}
	panic("invalid component")
}

///////////////////////////////////////////////////////////////////////////////
