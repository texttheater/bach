package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/ast"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
)

///////////////////////////////////////////////////////////////////////////////

type Definition struct {
	Pos         lexer.Position
	InputType   *Type        `"for" @@`
	NameParlist *NameParlist `"def" ( @@`
	Name        *string      `      | @Name)`
	OutputType  *Type        `@@`
	Body        *Composition `"as" @@ "ok"`
}

func (g *Definition) Ast() ast.Expression {
	var name = g.Name
	var params []*functions.Parameter
	if g.NameParlist != nil {
		name = &g.NameParlist.NameLpar.Name
		params = g.NameParlist.Ast()
	}
	return &ast.DefinitionExpression{
		Pos:        g.Pos,
		InputType:  g.InputType.Ast(),
		Name:       *name,
		Params:     params,
		OutputType: g.OutputType.Ast(),
		Body:       g.Body.Ast(),
	}
}

///////////////////////////////////////////////////////////////////////////////

type NameParlist struct {
	Pos      lexer.Position
	NameLpar *NameLpar `@NameLpar`
	Param    *Param    `@@`
	Params   []*Param  `{ "," @@ } ")"`
}

func (g *NameParlist) Ast() []*functions.Parameter {
	params := make([]*functions.Parameter, 0, 1+len(g.Params))
	params = append(params, g.Param.Ast())
	for _, param := range g.Params {
		params = append(params, param.Ast())
	}
	return params
}

///////////////////////////////////////////////////////////////////////////////

type Param struct {
	Pos         lexer.Position
	Name1       *string      `( @Name`
	InputType   *Type        `| "for" @@`
	NameParlist *NameParlist `  ( @@`
	Name2       *string      `  | @Name) )`
	OutputType  *Type        `@@`
}

func (g *Param) Ast() *functions.Parameter {
	var inputType types.Type
	if g.InputType != nil {
		inputType = g.InputType.Ast()
	} else {
		inputType = &types.AnyType{}
	}
	var name string
	var params []*functions.Parameter
	if g.Name1 != nil {
		name = *g.Name1
	} else if g.NameParlist != nil {
		name = g.NameParlist.NameLpar.Name
		params = make([]*functions.Parameter, 0, len(g.NameParlist.Params)+1)
		params = append(params, g.NameParlist.Param.Ast())
		for _, param := range g.NameParlist.Params {
			params = append(params, param.Ast())
		}
	} else {
		name = *g.Name2
	}
	return &functions.Parameter{
		InputType:  inputType,
		Name:       name,
		Params:     params,
		OutputType: g.OutputType.Ast(),
	}
}

///////////////////////////////////////////////////////////////////////////////
