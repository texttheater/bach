package grammar

import (
	"strconv"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

///////////////////////////////////////////////////////////////////////////////

type Call struct {
	Pos         lexer.Position
	Op1Num      *Op1Num      `  @Op1Num`
	Op2Num      *Op2Num      `| @Op2Num`
	Op1Name     *Op1Name     `| @Op1Name`
	Op2Name     *Op2Name     `| @Op2Name`
	NameArglist *NameArglist `| @@`
	Name        *string      `| ( @Prop | @Op1 | @Op2 )`
}

func (g *Call) Ast() (functions.Expression, error) {
	if g.Op1Num != nil {
		return g.Op1Num.Ast(), nil
	}
	if g.Op2Num != nil {
		return g.Op2Num.Ast(), nil
	}
	if g.Op1Name != nil {
		return g.Op1Name.Ast(), nil
	}
	if g.Op2Name != nil {
		return g.Op2Name.Ast(), nil
	}
	if g.NameArglist != nil {
		return g.NameArglist.Ast()
	}
	if g.Name != nil {
		return &functions.CallExpression{g.Pos, *g.Name, []functions.Expression{}}, nil
	}
	panic("invalid call")
}

///////////////////////////////////////////////////////////////////////////////

type Op1Num struct {
	Pos lexer.Position
	Op  string
	Num float64
}

func (g *Op1Num) Capture(values []string) error {
	g.Op = string(values[0][:1])
	f, err := strconv.ParseFloat(values[0][1:], 64)
	if err != nil {
		return err
	}
	g.Num = f
	return nil
}

func (g *Op1Num) Ast() functions.Expression {
	return &functions.CallExpression{
		Pos:  g.Pos,
		Name: g.Op,
		Args: []functions.Expression{
			&functions.ConstantExpression{
				Pos:   g.Pos,
				Type:  types.NumType{},
				Value: values.NumValue(g.Num),
			},
		},
	}
}

///////////////////////////////////////////////////////////////////////////////

type Op2Num struct {
	Pos lexer.Position
	Op  string
	Num float64
}

func (g *Op2Num) Capture(values []string) error {
	g.Op = string(values[0][:2])
	f, err := strconv.ParseFloat(values[0][2:], 64)
	if err != nil {
		return err
	}
	g.Num = f
	return nil
}

func (g *Op2Num) Ast() functions.Expression {
	return &functions.CallExpression{
		Pos:  g.Pos,
		Name: g.Op,
		Args: []functions.Expression{
			&functions.ConstantExpression{
				Pos:   g.Pos,
				Type:  types.NumType{},
				Value: values.NumValue(g.Num),
			},
		},
	}
}

///////////////////////////////////////////////////////////////////////////////

type Op1Name struct {
	Pos  lexer.Position
	Op   string
	Name string
}

func (g *Op1Name) Capture(values []string) error {
	g.Op = string(values[0][:1])
	g.Name = values[0][1:]
	return nil
}

func (g *Op1Name) Ast() functions.Expression {
	return &functions.CallExpression{g.Pos, g.Op, []functions.Expression{&functions.CallExpression{g.Pos, g.Name, []functions.Expression{}}}}
}

///////////////////////////////////////////////////////////////////////////////

type Op2Name struct {
	Pos  lexer.Position
	Op   string
	Name string
}

func (g *Op2Name) Capture(values []string) error {
	g.Op = string(values[0][:2])
	g.Name = values[0][2:]
	return nil
}

func (g *Op2Name) Ast() functions.Expression {
	return &functions.CallExpression{g.Pos, g.Op, []functions.Expression{&functions.CallExpression{g.Pos, g.Name, []functions.Expression{}}}}
}

///////////////////////////////////////////////////////////////////////////////

type NameArglist struct {
	Pos      lexer.Position
	NameLpar *NameLpar      `@NameLpar`
	Arg      *Composition   `@@`
	Args     []*Composition `( "," @@ )* ")"`
}

func (g *NameArglist) Ast() (functions.Expression, error) {
	args := make([]functions.Expression, len(g.Args)+1)
	var err error
	args[0], err = g.Arg.Ast()
	if err != nil {
		return nil, err
	}
	for i, Arg := range g.Args {
		args[i+1], err = Arg.Ast()
		if err != nil {
			return nil, err
		}
	}
	return &functions.CallExpression{g.Pos, g.NameLpar.Name, args}, nil
}

///////////////////////////////////////////////////////////////////////////////

type NameLpar struct {
	Pos  lexer.Position
	Name string
}

func (g *NameLpar) Capture(values []string) error {
	g.Name = values[0][:len(values[0])-1]
	return nil
}

///////////////////////////////////////////////////////////////////////////////
