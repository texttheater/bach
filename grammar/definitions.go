package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/types"
)

type Definition struct {
	Pos        lexer.Position
	InputType  *Type        `"for" @@`
	Name       string       `"def" ( ( @Prop | @Op1 | @Op2 )`
	NameLpar   *NameLpar    `      | @NameLpar`
	ParamName  *string      `        ( @Prop | @Op1 | @Op2 )`
	Parameter  *Parameter   `        @@`
	ParamNames []string     `        ( "," ( @Prop | @Op1 | @Op2 )`
	Params     []*Parameter `          @@ )* ")" )`
	OutputType *Type        `@@`
	Body       *Composition `"as" @@ "ok"`
}

func (g *Definition) Ast() (expressions.Expression, error) {
	inputType := g.InputType.Ast()
	var name = g.Name
	var params []*shapes.Parameter
	var paramNames []string
	if g.NameLpar != nil {
		name = g.NameLpar.Name
		params = make([]*shapes.Parameter, len(g.Params)+1)
		paramNames = make([]string, len(g.ParamNames)+1)
		param, err := g.Parameter.Ast()
		if err != nil {
			return nil, err
		}
		params[0] = param
		paramNames[0] = *g.ParamName
		for i, param := range g.Params {
			params[i+1], err = param.Ast()
			if err != nil {
				return nil, err
			}
			paramNames[i+1] = g.ParamNames[i+1]
		}
	}
	outputType := g.OutputType.Ast()
	body, err := g.Body.Ast()
	if err != nil {
		return nil, err
	}
	return &expressions.DefinitionExpression{
		Pos:        g.Pos,
		InputType:  inputType,
		Name:       name,
		Params:     params,
		ParamNames: paramNames,
		OutputType: outputType,
		Body:       body,
	}, nil
}

type Parameter struct {
	Pos        lexer.Position
	InputType  *Type        `( "for" @@`
	Parameter  *Parameter   `  ( "(" @@`
	Params     []*Parameter `    ( "," @@ )* ")" )? )?`
	OutputType *Type        `@@`
}

func (g *Parameter) Ast() (*shapes.Parameter, error) {
	var inputType types.Type
	if g.InputType != nil {
		inputType = g.InputType.Ast()
	} else {
		inputType = types.AnyType{}
	}
	var params []*shapes.Parameter
	if g.Parameter != nil {
		params = make([]*shapes.Parameter, len(g.Params)+1)
		var err error
		params[0], err = g.Parameter.Ast()
		if err != nil {
			return nil, err
		}
		for i, param := range g.Params {
			params[i+1], err = param.Ast()
			if err != nil {
				return nil, err
			}
		}
	}
	outputType := g.OutputType.Ast()
	return &shapes.Parameter{
		InputType:  inputType,
		Params:     params,
		OutputType: outputType,
	}, nil
}
