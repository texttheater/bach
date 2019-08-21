package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
)

type Filter struct {
	Pos               lexer.Position
	Pattern           *Pattern     `( "eachis" @@`
	Guard             *Composition `  ( "with" @@ )?`
	Condition         *Composition `| "eachif" @@ )`
	Consequent        *Composition `( "then" ( @@ | "drop" )` // long form
	LongAlternatives  []*FLongAlt  `  ( @@ )*`
	Alternative       *Composition `  ( "else" ( @@ | "drop" ) )?` // short form
	ShortAlternatives []*FShortAlt `| ( @@ )* )`
}

type FLongAlt struct {
	Pos        lexer.Position
	Pattern    *Pattern     `( "elis" @@`
	Guard      *Composition `  ( "with" @@ )?`
	Condition  *Composition `| "elif" @@ )`
	Consequent *Composition `"then" @@`
}

type FShortAlt struct {
	Pos       lexer.Position
	Pattern   *Pattern     `( "elis" @@`
	Guard     *Composition `  ( "with" @@ )?`
	Condition *Composition `| "elif" @@ )`
}

func (g *Filter) Ast() (functions.Expression, error) {
	if g.Consequent != nil { // long form
		panic("not implemented yet")
	} else { // short form
		panic("not implemented yet")
	}
}
