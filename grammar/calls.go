package grammar

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type Call struct {
	Pos         lexer.Position
	Op1Num      *Op1Num      `  @@`
	Op2Num      *Op2Num      `| @@`
	Op1Lid     *Op1Lid     `| @@`
	Op2Lid     *Op2Lid     `| @@`
	NameRegexp  *NameRegexp  `| @@`
	NameArglist *NameArglist `| @@`
	Name        *string      `| ( @Lid | @Op1 | @Op2 )`
}

func (g *Call) Ast() (functions.Expression, error) {
	if g.Op1Num != nil {
		return g.Op1Num.Ast()
	}
	if g.Op2Num != nil {
		return g.Op2Num.Ast()
	}
	if g.Op1Lid != nil {
		return g.Op1Lid.Ast()
	}
	if g.Op2Lid != nil {
		return g.Op2Lid.Ast()
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
	numPos := g.Pos
	numPos.Column += 1
	return &functions.CallExpression{
		Pos:  g.Pos,
		Name: op,
		Args: []functions.Expression{
			&functions.ConstantExpression{
				Pos:   numPos,
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
	numPos := g.Pos
	numPos.Column += 2
	return &functions.CallExpression{
		Pos:  g.Pos,
		Name: op,
		Args: []functions.Expression{
			&functions.ConstantExpression{
				Pos:   numPos,
				Type:  types.NumType{},
				Value: values.NumValue(num),
			},
		},
	}, nil
}

type Op1Lid struct {
	Pos     lexer.Position
	Op1Lid *string `@Op1Lid`
}

func (g *Op1Lid) Ast() (functions.Expression, error) {
	op1name := *g.Op1Lid
	op := op1name[:1]
	name := op1name[1:]
	namePos := g.Pos
	namePos.Column += 1
	return &functions.CallExpression{
		Pos:  g.Pos,
		Name: op,
		Args: []functions.Expression{
			&functions.CallExpression{
				Pos:  namePos,
				Name: name,
				Args: nil,
			},
		},
	}, nil
}

type Op2Lid struct {
	Pos     lexer.Position
	Op2Lid *string `@Op2Lid`
}

func (g *Op2Lid) Ast() (functions.Expression, error) {
	op2name := *g.Op2Lid
	op := op2name[:2]
	name := op2name[2:]
	namePos := g.Pos
	namePos.Column += 2
	return &functions.CallExpression{
		Pos:  g.Pos,
		Name: op,
		Args: []functions.Expression{
			&functions.CallExpression{
				Pos:  namePos,
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
	regexpPos := g.Pos
	regexpPos.Column += len(name)
	regexp, err := regexp.Compile(regexpString)
	if err != nil {
		return nil, errors.E(
			errors.Pos(regexpPos),
			errors.Code(errors.BadRegexp),
			errors.Message(err.Error()),
		)
	}
	return &functions.CallExpression{
		Pos:  g.Pos,
		Name: name,
		Args: []functions.Expression{
			&functions.RegexpExpression{
				Pos:    regexpPos,
				Regexp: regexp,
			},
		},
	}, nil
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
