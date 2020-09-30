package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type TemplateLiteral struct {
	Pos       lexer.Position
	Fragments []*Fragment "\"`\" @@* \"`\""
}

func (g *TemplateLiteral) Ast() (expressions.Expression, error) {
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

type Fragment struct {
	Pos         lexer.Position
	Composition *Composition `( "{" @@ "}"`
	Text        string       `| @Char )`
}

func (g *Fragment) Ast() (expressions.Expression, error) {
	if g.Composition != nil {
		return g.Composition.Ast()
	}
	return &expressions.ConstantExpression{
		Pos:   g.Pos,
		Type:  types.StrType{},
		Value: states.StrValue(g.Text), // TODO allow escapes?
	}, nil
}
