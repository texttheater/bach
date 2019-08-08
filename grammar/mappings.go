package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
)

type Mapping struct {
	Pos       lexer.Position    `"each"`
	Body      *Composition      `( @@`
	Pattern   *Pattern          `| ( @@`
	Guard     *Composition      `    ( "with" @@)?`
	Condition *Composition      `  | "if" @@)`
	LongTail  *LongMappingTail  `  ( @@`
	ShortTail *ShortMappingTail `  | @@) ) "all"`
}

type ShortMappingTail struct {
	Pos          lexer.Position
	Alternatives []*ShortMappingAlternative `( @@ )*`
}

type LongMappingTail struct {
	Pos          lexer.Position
	Consequent   *Composition              `"then" ( @@ | "drop")`
	Alternatives []*LongMappingAlternative `( @@ )*`
}

type ShortMappingAlternative struct {
	Pos       lexer.Position
	Pattern   *Pattern     `( "elis" @@`
	Guard     *Composition `  ( "with" @@ )?`
	Condition *Composition `| "elif" @@ )`
}

type LongMappingAlternative struct {
	Pos        lexer.Position
	Pattern    *Pattern     `( "elis" @@`
	Guard      *Composition `  ( "with" @@ )?`
	Condition  *Composition `| "elif" @@ )`
	Consequent *Composition `"then" @@`
}

func (g *Mapping) Ast() (functions.Expression, error) {
	body, err := g.Body.Ast() // TODO handle other case
	if err != nil {
		return nil, err
	}
	return &functions.MappingExpression{
		Pos:  g.Pos,
		Body: body,
	}, nil
}
