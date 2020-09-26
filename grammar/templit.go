package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type TemplateLiteral struct {
	Pos       lexer.Position
	Fragments []*Fragment "\"`\" @@* \"`\""
}

func (g *TemplateLiteral) Ast() (functions.Expression, error) {
	pieces := make([]functions.Expression, len(g.Fragments))
	for i, fragment := range g.Fragments {
		piece, err := fragment.Ast()
		if err != nil {
			return nil, err
		}
		pieces[i] = piece
	}
	return &functions.TemplateLiteralExpression{
		Pos:    g.Pos,
		Pieces: pieces,
	}, nil
}

type Fragment struct {
	Pos         lexer.Position
	Composition *Composition `( "{" @@ "}"`
	Text        string       `| @Char )`
}

func (g *Fragment) Ast() (functions.Expression, error) {
	if g.Composition != nil {
		return g.Composition.Ast()
	}
	return &functions.ConstantExpression{
		Pos:   g.Pos,
		Type:  types.StrType{},
		Value: states.StrValue(g.Text), // TODO allow escapes?
	}, nil
}
