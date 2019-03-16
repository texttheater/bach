package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
)

///////////////////////////////////////////////////////////////////////////////

type Definition struct {
	Pos         lexer.Position
	InputType   *Type        `"for" @@`
	NameParlist *NameParlist `"def" ( @@`
	Name        *string      `      | ( @Prop | @Op1 | @Op2 ) )`
	OutputType  *Type        `@@`
	Body        *Composition `"as" @@ "ok"`
}

func (g *Definition) Ast() expressions.Expression {
	var name = g.Name
	var params []*functions.Parameter
	if g.NameParlist != nil {
		name = &g.NameParlist.NameLpar.Name
		params = g.NameParlist.Ast()
	}
	return &expressions.DefinitionExpression{
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
	params := make([]*functions.Parameter, 1+len(g.Params))
	params[0] = g.Param.Ast()
	for i, param := range g.Params {
		params[i+1] = param.Ast()
	}
	return params
}

///////////////////////////////////////////////////////////////////////////////

type Param struct {
	Pos         lexer.Position
	Name1       *string      `( ( @Prop | @Op1 | @Op2 )`
	InputType   *Type        `| "for" @@`
	NameParlist *NameParlist `  ( @@`
	Name2       *string      `  | ( @Prop | @Op1 | @Op2 ) ) )`
	OutputType  *Type        `@@`
}

func (g *Param) Ast() *functions.Parameter {
	var inputType types.Type
	if g.InputType != nil {
		inputType = g.InputType.Ast()
	} else {
		inputType = types.AnyType
	}
	var name string
	var params []*functions.Parameter
	if g.Name1 != nil {
		name = *g.Name1
	} else if g.NameParlist != nil {
		name = g.NameParlist.NameLpar.Name
		params = make([]*functions.Parameter, len(g.NameParlist.Params)+1)
		params[0] = g.NameParlist.Param.Ast()
		for i, param := range g.NameParlist.Params {
			params[i+1] = param.Ast()
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
