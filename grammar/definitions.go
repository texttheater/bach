package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
)

///////////////////////////////////////////////////////////////////////////////

type Definition struct {
	Pos         lexer.Position
	InputType   *Type        `"for" @@`
	NameParlist *NameParlist `"def" ( @@`
	Name        *string      `      | @Name)`
	OutputType  *Type
	Body        *Composition `"as" @@ "ok"`
}

///////////////////////////////////////////////////////////////////////////////

type NameParlist struct {
	Pos      lexer.Position
	NameLpar *NameLpar `@NameLpar`
	Param    *Param    `@@`
	Params   []*Param  `{ "," @@ } ")"`
}

///////////////////////////////////////////////////////////////////////////////

type Param struct {
	Pos         lexer.Position
	InputType   *Type        `"for" @@`
	NameParlist *NameParlist `"def" ( @@`
	Name        *string      `      | @Name)`
	OutputType  *Type        `@@`
}

func (g *Param) Ast() *functions.Param {
	var name string
	var params []*functions.Param
	if g.NameParlist != nil {
		name = g.NameParlist.NameLpar.Name
		params = make([]*functions.Param, 0, len(g.NameParlist.Params)+1)
		params = append(params, g.NameParlist.Param.Ast())
		for _, param := range g.NameParlist.Params {
			params = append(params, param.Ast())
		}
	} else {
		name = *g.Name
		params = nil
	}
	return &functions.Param{
		InputType:  g.InputType.Ast(),
		Name:       name,
		Params:     params,
		OutputType: g.OutputType.Ast(),
	}
}

///////////////////////////////////////////////////////////////////////////////

type Type struct {
	Pos lexer.Position `"Num"`
}

func (g *Type) Ast() types.Type {
	return &types.NumberType{}
}

///////////////////////////////////////////////////////////////////////////////
