package grammar

import (
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

///////////////////////////////////////////////////////////////////////////////

var LexerDefinition = lexer.Must(lexer.Regexp(
	`([\s]+)` +
		`|(?P<Num>(?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+))` +
		`|(?P<Str>"(?:\\.|[^"])*")` +
		`|(?P<Op1Num>[+\-*/%<>](?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+))` +
		`|(?P<Op2Num>(?:==|<=|>=)(?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+))` +
		`|(?P<Op1Name>[+\-*/%<>](?:[\p{L}_][\p{L}_0-9]*))` +
		`|(?P<Op2Name>(?:==|<=|>=)(?:[\p{L}_][\p{L}_0-9]*))` +
		`|(?P<Assignment>=(?:[+\-*/%<>=]|==|<=|>=|[\p{L}_][\p{L}_0-9]*))` +
		`|(?P<NameLpar>(?:[+\-*/%<>=]|==|<=|>=|[\p{L}_][\p{L}_0-9]*)\()` +
		`|(?P<TypeKeywordLangle>(?:Void|Null|Bool|Num|Str|Seq|Arr|Obj|Any)<)` +
		`|(?P<Prop>[\p{L}_][\p{L}_0-9]*)` +
		`|(?P<Op1>[+\-*/%<>=])` +
		`|(?P<Op2>==|<=|>=)` +
		`|(?P<Comma>,)` +
		`|(?P<Rpar>\))` +
		`|(?P<Lbrack>\[)` +
		`|(?P<Rbrack>])` +
		`|(?P<Lbrace>{)` +
		`|(?P<Rbrace>})` +
		`|(?P<Colon>:)` +
		`|(?P<Pipe>\|)` +
		// the following will be scanned as Name, but mapped to the
		// appropriate token types by name2keyword (see below)
		`|(?P<Keyword>for|def|as|ok|if|then|elif|else)` +
		`|(?P<TypeKeyword>Void|Null|Bool|Num|Str|Seq|Arr|Obj|Any)`,
))

///////////////////////////////////////////////////////////////////////////////

func Parse(input string) (expressions.Expression, error) {
	parser, err := participle.Build(
		&Composition{},
		participle.Lexer(LexerDefinition),
		participle.Unquote("Str"),
		participle.Map(ToKeyword, "Prop"),
		participle.UseLookahead(0),
	)
	if err != nil {
		if lexerError, ok := err.(*lexer.Error); ok {
			return nil, errors.E(
				errors.Kind(errors.Syntax),
				errors.Pos(lexerError.Pos),
				errors.Message(lexerError.Message),
			)
		}
		return nil, err
	}
	composition := &Composition{}
	err = parser.ParseString(input, composition)
	if err != nil {
		if lexerError, ok := err.(*lexer.Error); ok {
			return nil, errors.E(
				errors.Kind(errors.Syntax),
				errors.Pos(lexerError.Pos),
				errors.Message(lexerError.Message),
			)
		}
		return nil, err
	}
	return composition.Ast(), nil
}

func ToKeyword(t lexer.Token) (lexer.Token, error) {
	if isTypeKeyword(t.Value) {
		t.Type = LexerDefinition.Symbols()["Type"]
	}
	if isKeyword(t.Value) {
		t.Type = LexerDefinition.Symbols()["Keyword"]
	}
	return t, nil
}

func isTypeKeyword(name string) bool {
	return name == "Void" || name == "Null" || name == "Bool" ||
		name == "Num" || name == "Str" || name == "Seq" ||
		name == "Arr" || name == "Any"
}

func isKeyword(name string) bool {
	return name == "for" || name == "def" || name == "as" ||
		name == "ok" || name == "if" || name == "then" ||
		name == "elif" || name == "else" || name == "each" ||
		name == "all"
}

///////////////////////////////////////////////////////////////////////////////

type Composition struct {
	Pos        lexer.Position
	Component  *Component   `@@`
	Components []*Component `{ @@ }`
}

func (g *Composition) Ast() expressions.Expression {
	pos := g.Component.Pos
	e := g.Component.Ast()
	for _, comp := range g.Components {
		e = &expressions.CompositionExpression{pos, e, comp.Ast()}
	}
	return e
}

///////////////////////////////////////////////////////////////////////////////

type Component struct {
	Pos         lexer.Position
	Num         *float64     `  @Num`
	Str         *string      `| @Str`
	Array       *Array       `| @@`
	Object      *Object      `| @@`
	Call        *Call        `| @@`
	Assignment  *Assignment  `| @Assignment`
	Definition  *Definition  `| @@`
	Conditional *Conditional `| @@`
	Mapping     *Mapping     `| @@`
}

func (g *Component) Ast() expressions.Expression {
	if g.Num != nil {
		return &expressions.ConstantExpression{
			Pos:   g.Pos,
			Type:  types.NumType,
			Value: values.NumValue(*g.Num),
		}
	}
	if g.Str != nil {
		return &expressions.ConstantExpression{
			Pos:   g.Pos,
			Type:  types.StrType,
			Value: values.StrValue(*g.Str),
		}
	}
	if g.Array != nil {
		return g.Array.Ast()
	}
	if g.Object != nil {
		return g.Object.Ast()
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
	if g.Mapping != nil {
		return g.Mapping.Ast()
	}
	panic("invalid component")
}

///////////////////////////////////////////////////////////////////////////////
