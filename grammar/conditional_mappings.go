package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
)

type ConditionalMapping struct {
	Pos               lexer.Position
	Pattern           *Pattern      `( "eachis" @@`
	Guard             *Composition  `  ( "with" @@ )?`
	Condition         *Composition  `| "eachif" @@ )`
	Consequent        *Composition  `( "then" ( @@ | "drop" )`
	LongAlternatives  []*CMLongAlt  `  ( @@ )*`
	Alternative       *Composition  `  ( "else" ( @@ | "drop" ) )?`
	ShortAlternatives []*CMShortAlt `| ( @@ )*`
}

type CMLongAlt struct {
	Pos        lexer.Position
	Pattern    *Pattern     `( "elis" @@`
	Guard      *Composition `  ( "with" @@ )?`
	Condition  *Composition `| "elif" @@ )`
	Consequent *Composition `"then" @@`
}

type CMShortAlt struct {
	Pos       lexer.Position
	Pattern   *Pattern     `( "elis" @@`
	Guard     *Composition `  ( "with" @@ )?`
	Condition *Composition `| "elif" @@ )`
}

func (g *ConditionalMapping) Ast() (functions.Expression, error) {
	panic("not implemented yet")
}
