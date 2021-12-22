package grammar

import (
	"strconv"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type StrLiteral struct {
	Pos       lexer.Position
	Fragments []*Fragment `"\"" @@* "\""`
}

func (g *StrLiteral) Ast() (expressions.Expression, error) {
	pieces := make([]expressions.Expression, len(g.Fragments))
	for i, fragment := range g.Fragments {
		piece, err := fragment.Ast()
		if err != nil {
			return nil, err
		}
		pieces[i] = piece
	}
	return &expressions.TemplateLiteralExpression{
		Pos:    g.Pos,
		Pieces: pieces,
	}, nil
}

func (g *StrLiteral) StaticStr() (string, bool, error) {
	if len(g.Fragments) != 1 {
		return "", false, nil
	}
	fragment := g.Fragments[0]
	if fragment.Dbrace != nil {
		if len(*fragment.Dbrace) == 1 {
			return "", false, errors.SyntaxError(
				errors.Code(errors.Syntax),
				errors.Pos(g.Pos),
				errors.Message("Use a double brace }} for a literal brace }"),
			)
		}
		return (*fragment.Dbrace)[:1], true, nil
	}
	if fragment.Text != nil {
		str, err := strconv.Unquote("\"" + *fragment.Text + "\"")
		if err != nil {
			return "", false, errors.SyntaxError(
				errors.Code(errors.Syntax),
				errors.Pos(g.Pos),
				errors.Message(err.Error()),
			)
		}
		return str, true, nil
	}
	return "", false, nil
}

type Fragment struct {
	Pos         lexer.Position
	Composition *Composition `( "{" @@ "}"`
	Dbrace      *string      `| @Dbrace`
	Text        *string      `| @Char )`
}

func (g *Fragment) Ast() (expressions.Expression, error) {
	if g.Composition != nil {
		return g.Composition.Ast()
	}
	var str string
	var err error
	if g.Dbrace != nil {
		if len(*g.Dbrace) == 1 {
			return nil, errors.SyntaxError(
				errors.Code(errors.Syntax),
				errors.Pos(g.Pos),
				errors.Message("Use a double brace }} for a literal brace }"),
			)
		}
		str = (*g.Dbrace)[:1]
	} else {
		str, err = strconv.Unquote("\"" + *g.Text + "\"")
		if err != nil {
			return nil, errors.SyntaxError(
				errors.Code(errors.Syntax),
				errors.Pos(g.Pos),
				errors.Message(err.Error()),
			)
		}
	}
	return &expressions.ConstantExpression{
		Pos:   g.Pos,
		Type:  types.Str{},
		Value: states.StrValue(str),
	}, nil
}
