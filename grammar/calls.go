package grammar

import (
	"strconv"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/ast"
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
	Name        *string      `| @Name`
}

func (g *Call) Ast() ast.Expression {
	if g.Op1Num != nil {
		return g.Op1Num.Ast()
	}
	if g.Op2Num != nil {
		return g.Op2Num.Ast()
	}
	if g.Op1Name != nil {
		return g.Op1Name.Ast()
	}
	if g.Op2Name != nil {
		return g.Op2Name.Ast()
	}
	if g.NameArglist != nil {
		return g.NameArglist.Ast()
	}
	if g.Name != nil {
		return &ast.CallExpression{g.Pos, *g.Name, []ast.Expression{}}
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

func (g *Op1Num) Ast() ast.Expression {
	return &ast.CallExpression{
		Pos:  g.Pos,
		Name: g.Op,
		Args: []ast.Expression{
			&ast.ConstantExpression{
				Pos:   g.Pos,
				Type:  &types.NumType{},
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

func (g *Op2Num) Ast() ast.Expression {
	return &ast.CallExpression{
		Pos:  g.Pos,
		Name: g.Op,
		Args: []ast.Expression{
			&ast.ConstantExpression{
				Pos:   g.Pos,
				Type:  &types.NumType{},
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

func (g *Op1Name) Ast() ast.Expression {
	return &ast.CallExpression{g.Pos, g.Op, []ast.Expression{&ast.CallExpression{g.Pos, g.Name, []ast.Expression{}}}}
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

func (g *Op2Name) Ast() ast.Expression {
	return &ast.CallExpression{g.Pos, g.Op, []ast.Expression{&ast.CallExpression{g.Pos, g.Name, []ast.Expression{}}}}
}

///////////////////////////////////////////////////////////////////////////////

type NameArglist struct {
	Pos      lexer.Position
	NameLpar *NameLpar      `@NameLpar`
	Arg      *Composition   `@@`
	Args     []*Composition `{ "," @@ } ")"`
}

func (g *NameArglist) Ast() ast.Expression {
	args := make([]ast.Expression, len(g.Args)+1)
	args[0] = g.Arg.Ast()
	for i, Arg := range g.Args {
		args[i+1] = Arg.Ast()
	}
	return &ast.CallExpression{g.Pos, g.NameLpar.Name, args}
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
