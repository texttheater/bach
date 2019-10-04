package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
)

type Reject struct {
	Pos lexer.Position `"reject"`
}

func (g *Reject) Ast() (functions.Expression, error) {
	return &functions.RejectExpression{
		Pos: g.Pos,
	}, nil
}
