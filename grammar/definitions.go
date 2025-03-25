package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/params"
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
