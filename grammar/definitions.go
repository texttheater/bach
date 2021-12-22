package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/parameters"
	"github.com/texttheater/bach/types"
)

type Definition struct {
	Pos        lexer.Position
	InputType  *TypeTemplate     `"for" @@`
	Name       *string           `"def" ( ( @Lid | @Op1 | @Op2 )`
	NameLpar   *string           `      | @NameLpar`
	Parameter  *NamedParameter   `        @@`
	Params     []*NamedParameter `        ( "," @@ )* ")" )`
	OutputType *Type             `@@`
	Body       *Composition      `"as" @@ "ok"`
}

func (g *Definition) Ast() (expressions.Expression, error) {
	inputType := g.InputType.Ast()
	var name string
	var params []*parameters.Parameter
	var paramNames []string
	if g.Name != nil {
		name = *g.Name
	} else {
		nameLpar := *g.NameLpar
		name = nameLpar[:len(nameLpar)-1]
		params = make([]*parameters.Parameter, len(g.Params)+1)
		paramNames = make([]string, len(g.Params)+1)
		param, paramName, err := g.Parameter.Ast()
		if err != nil {
			return nil, err
		}
		params[0] = param
		paramNames[0] = paramName
		for i, param := range g.Params {
			params[i+1], paramNames[i+1], err = param.Ast()
			if err != nil {
				return nil, err
			}
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

type NamedParameter struct {
	Pos   lexer.Position
	Long  *NamedParameterLongForm  `( @@`
	Short *NamedParameterShortForm `| @@ )`
}

func (g *NamedParameter) Ast() (*parameters.Parameter, string, error) {
	if g.Long != nil {
		return g.Long.Ast()
	}
	return g.Short.Ast()
}

type NamedParameterLongForm struct {
	Pos        lexer.Position
	InputType  *TypeTemplate `"for" @@`
	Name       *string       `( ( @Lid | @Op1 | @Op2 )`
	NameLpar   *string       `| @NameLpar`
	Parameter  *Parameter    `  @@`
	Params     []*Parameter  `  ( "," @@ )* ")" )`
	OutputType *TypeTemplate `@@`
}

func (g *NamedParameterLongForm) Ast() (*parameters.Parameter, string, error) {
	inputType := g.InputType.Ast()
	var name string
	var params []*parameters.Parameter
	if g.Name != nil {
		name = *g.Name
	} else {
		name = (*g.NameLpar)[:len(*g.NameLpar)-1]
		params = make([]*parameters.Parameter, len(g.Params)+1)
		var err error
		params[0], err = g.Parameter.Ast()
		if err != nil {
			return nil, "", err
		}
		for i := range g.Params {
			params[i+1], err = g.Params[i].Ast()
			if err != nil {
				return nil, "", err
			}
		}
	}
	outputType := g.OutputType.Ast()
	return &parameters.Parameter{
		InputType:  inputType,
		Params:     params,
		OutputType: outputType,
	}, name, nil
}

type NamedParameterShortForm struct {
	Pos        lexer.Position
	Name       string        `( @Lid | @Op1 | @Op2)`
	OutputType *TypeTemplate `@@`
}

func (g *NamedParameterShortForm) Ast() (*parameters.Parameter, string, error) {
	outputType := g.OutputType.Ast()
	return &parameters.Parameter{
		InputType:  types.Any{},
		OutputType: outputType,
	}, g.Name, nil
}

type Parameter struct {
	Pos        lexer.Position
	InputType  *TypeTemplate `( "for" @@`
	Parameter  *Parameter    `  ( "(" @@`
	Params     []*Parameter  `    ( "," @@ )* ")" )? )?`
	OutputType *TypeTemplate `@@`
}

func (g *Parameter) Ast() (*parameters.Parameter, error) {
	var inputType types.Type
	if g.InputType != nil {
		inputType = g.InputType.Ast()
	} else {
		inputType = types.Any{}
	}
	var params []*parameters.Parameter
	if g.Parameter != nil {
		params = make([]*parameters.Parameter, len(g.Params)+1)
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
	return &parameters.Parameter{
		InputType:  inputType,
		Params:     params,
		OutputType: outputType,
	}, nil
}
