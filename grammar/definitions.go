package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/types"
)

type Definition struct {
	Pos        lexer.Position
	InputType  *TypeTemplate `"for" @@`
	Name       *string       `"def" ( ( @Lid | @Op1 | @Op2 )`
	NameLpar   *string       `      | @NameLpar`
	Param      *NamedParam   `        @@`
	Params     []*NamedParam `        ( "," @@ )* ")" )`
	OutputType *Type         `@@`
	Body       *Composition  `"as" @@ "ok"`
}

func (g *Definition) Ast() (expressions.Expression, error) {
	inputType := g.InputType.Ast()
	var name string
	var pars []*params.Param
	var parNames []string
	if g.Name != nil {
		name = *g.Name
	} else {
		nameLpar := *g.NameLpar
		name = nameLpar[:len(nameLpar)-1]
		pars = make([]*params.Param, len(g.Params)+1)
		parNames = make([]string, len(g.Params)+1)
		par, parName, err := g.Param.Ast()
		if err != nil {
			return nil, err
		}
		pars[0] = par
		parNames[0] = parName
		for i, par := range g.Params {
			pars[i+1], parNames[i+1], err = par.Ast()
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
		Params:     pars,
		ParamNames: parNames,
		OutputType: outputType,
		Body:       body,
	}, nil
}

type NamedParam struct {
	Pos   lexer.Position
	Long  *NamedParamLongForm  `( @@`
	Short *NamedParamShortForm `| @@ )`
}

func (g *NamedParam) Ast() (*params.Param, string, error) {
	if g.Long != nil {
		return g.Long.Ast()
	}
	return g.Short.Ast()
}

type NamedParamLongForm struct {
	Pos        lexer.Position
	InputType  *TypeTemplate `"for" @@`
	Name       *string       `( ( @Lid | @Op1 | @Op2 )`
	NameLpar   *string       `| @NameLpar`
	Param      *Param        `  @@`
	Params     []*Param      `  ( "," @@ )* ")" )`
	OutputType *TypeTemplate `@@`
}

func (g *NamedParamLongForm) Ast() (*params.Param, string, error) {
	inputType := g.InputType.Ast()
	var name string
	var pars []*params.Param
	if g.Name != nil {
		name = *g.Name
	} else {
		name = (*g.NameLpar)[:len(*g.NameLpar)-1]
		pars = make([]*params.Param, len(g.Params)+1)
		var err error
		pars[0], err = g.Param.Ast()
		if err != nil {
			return nil, "", err
		}
		for i := range g.Params {
			pars[i+1], err = g.Params[i].Ast()
			if err != nil {
				return nil, "", err
			}
		}
	}
	outputType := g.OutputType.Ast()
	return &params.Param{
		InputType:  inputType,
		Params:     pars,
		OutputType: outputType,
	}, name, nil
}

type NamedParamShortForm struct {
	Pos        lexer.Position
	Name       string        `( @Lid | @Op1 | @Op2)`
	OutputType *TypeTemplate `@@`
}

func (g *NamedParamShortForm) Ast() (*params.Param, string, error) {
	outputType := g.OutputType.Ast()
	return &params.Param{
		InputType:  types.Any{},
		OutputType: outputType,
	}, g.Name, nil
}

type Param struct {
	Pos        lexer.Position
	InputType  *TypeTemplate `( "for" @@`
	Param      *Param        `  ( "(" @@`
	Params     []*Param      `    ( "," @@ )* ")" )? )?`
	OutputType *TypeTemplate `@@`
}

func (g *Param) Ast() (*params.Param, error) {
	var inputType types.Type
	if g.InputType != nil {
		inputType = g.InputType.Ast()
	} else {
		inputType = types.Any{}
	}
	var pars []*params.Param
	if g.Param != nil {
		pars = make([]*params.Param, len(g.Params)+1)
		var err error
		pars[0], err = g.Param.Ast()
		if err != nil {
			return nil, err
		}
		for i, param := range g.Params {
			pars[i+1], err = param.Ast()
			if err != nil {
				return nil, err
			}
		}
	}
	outputType := g.OutputType.Ast()
	return &params.Param{
		InputType:  inputType,
		Params:     pars,
		OutputType: outputType,
	}, nil
}
