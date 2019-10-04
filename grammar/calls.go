package grammar

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type Call struct {
	Pos         lexer.Position
	Op1Num      *Op1Num      `  @@`
	Op2Num      *Op2Num      `| @@`
	Op1Name     *Op1Name     `| @@`
	Op2Name     *Op2Name     `| @@`
	NameRegexp  *NameRegexp  `| @@`
	NameArglist *NameArglist `| @@`
	Name        *string      `| ( @Prop | @Op1 | @Op2 )`
}

func (g *Call) Ast() (functions.Expression, error) {
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
	if g.NameRegexp != nil {
		return g.NameRegexp.Ast()
	}
	if g.NameArglist != nil {
		return g.NameArglist.Ast()
	}
	if g.Name != nil {
		return &functions.CallExpression{
			Pos:  g.Pos,
			Name: *g.Name,
			Args: nil,
		}, nil
	}
	panic("invalid call")
}

type Op1Num struct {
	Pos    lexer.Position
	Op1Num *string `@Op1Num`
}

func (g *Op1Num) Ast() (functions.Expression, error) {
	op1num := *g.Op1Num
	op := op1num[:1]
	num, err := strconv.ParseFloat(op1num[1:], 64)
	if err != nil {
		return nil, err
	}
	return &functions.CallExpression{
		Pos:  g.Pos,
		Name: op,
		Args: []functions.Expression{
			&functions.ConstantExpression{
				Pos:   g.Pos,
				Type:  types.NumType{},
				Value: values.NumValue(num),
			},
		},
	}, nil
}

type Op2Num struct {
	Pos    lexer.Position
	Op2Num *string `@Op2Num`
}

func (g *Op2Num) Ast() (functions.Expression, error) {
	op2num := *g.Op2Num
	op := op2num[:2]
	num, err := strconv.ParseFloat(op2num[2:], 64)
	if err != nil {
		return nil, err
	}
	return &functions.CallExpression{
		Pos:  g.Pos,
		Name: op,
		Args: []functions.Expression{
			&functions.ConstantExpression{
				Pos:   g.Pos,
				Type:  types.NumType{},
				Value: values.NumValue(num),
			},
		},
	}, nil
}

type Op1Name struct {
	Pos     lexer.Position
	Op1Name *string `@Op1Name`
}

func (g *Op1Name) Ast() (functions.Expression, error) {
	op1name := *g.Op1Name
	op := op1name[:1]
	name := op1name[1:]
	return &functions.CallExpression{
		Pos:  g.Pos,
		Name: op,
		Args: []functions.Expression{
			&functions.CallExpression{
				Pos:  g.Pos,
				Name: name,
				Args: nil,
			},
		},
	}, nil
}

type Op2Name struct {
	Pos     lexer.Position
	Op2Name *string `@Op2Name`
}

func (g *Op2Name) Ast() (functions.Expression, error) {
	op2name := *g.Op2Name
	op := op2name[:2]
	name := op2name[2:]
	return &functions.CallExpression{
		Pos:  g.Pos,
		Name: op,
		Args: []functions.Expression{
			&functions.CallExpression{
				Pos:  g.Pos,
				Name: name,
				Args: nil,
			},
		},
	}, nil
}

type NameRegexp struct {
	Pos        lexer.Position
	NameRegexp *string `@NameRegexp`
}

func (g *NameRegexp) Ast() (functions.Expression, error) {
	nameRegexp := *g.NameRegexp
	i := strings.Index(nameRegexp, "~")
	name := nameRegexp[:i]
	regexpString := nameRegexp[i+1 : len(nameRegexp)-1]
	regexp, err := regexp.Compile(regexpString)
	if err != nil {
		return nil, err
	}
	regexpExpression := &functions.RegexpExpression{
		Pos:    g.Pos, // FIXME
		Regexp: regexp,
	}
	callExpression := &functions.CallExpression{
		Pos:  g.Pos,
		Name: name,
		Args: []functions.Expression{
			regexpExpression,
		},
	}
	return callExpression, nil
}

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

type NameLpar struct {
	Pos  lexer.Position
	Name string
}

func (g *NameLpar) Capture(values []string) error {
	g.Name = values[0][:len(values[0])-1]
	return nil
}
