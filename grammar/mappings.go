package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
)

type Mapping struct {
	Pos  lexer.Position
	Body *Composition `"each" @@ "all"`
}

func (g *Mapping) Ast() (functions.Expression, error) {
	body, err := g.Body.Ast()
	if err != nil {
		return nil, err
	}
	return &functions.MappingExpression{
		Pos:  g.Pos,
		Body: body,
	}, nil
}
