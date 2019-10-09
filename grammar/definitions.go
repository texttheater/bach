package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
)

type Definition struct {
	Pos        lexer.Position
	TypeParam  *string      `"for" ( ( @LangleLid | "<" @Lid )`
	TypeParams []string     `        ( "," @Lid )* ">" )?`
	InputType  *Type        `@@`
	Name       *string      `"def" ( ( @Lid | @Op1 | @Op2 )`
	NameLpar   *string      `      | @NameLpar`
	ParamName  *string      `        ( @Lid | @Op1 | @Op2 )`
	Parameter  *Parameter   `        @@`
	ParamNames []string     `        ( "," ( @Lid | @Op1 | @Op2 )`
	Params     []*Parameter `          @@ )* ")" )`
	OutputType *Type        `@@`
	Body       *Composition `"as" @@ "ok"`
}

func (g *Definition) Ast() (functions.Expression, error) {
	inputType := g.InputType.Ast()
	var name string
	var params []*functions.Parameter
	var paramNames []string
	if g.Name != nil {
		name = *g.Name
	} else {
		nameLpar := *g.NameLpar
		name = nameLpar[:len(nameLpar)-1]
		params = make([]*functions.Parameter, len(g.Params)+1)
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
	return &functions.DefinitionExpression{
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

func (g *Parameter) Ast() (*functions.Parameter, error) {
	var inputType types.Type
	if g.InputType != nil {
		inputType = g.InputType.Ast()
	} else {
		inputType = types.AnyType{}
	}
	var params []*functions.Parameter
	if g.Parameter != nil {
		params = make([]*functions.Parameter, len(g.Params)+1)
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
	return &functions.Parameter{
		InputType:  inputType,
		Params:     params,
		OutputType: outputType,
	}, nil
}
