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
	Op1Lid      *Op1Lid      `| @@`
	Op2Lid      *Op2Lid      `| @@`
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
	Op1Num string `@Op1Num`
}

func (g *Op1Num) Ast() (functions.Expression, error) {
	op := g.Op1Num[:1]
	num, err := strconv.ParseFloat(g.Op1Num[1:], 64)
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
	Op2Num string `@Op2Num`
}

func (g *Op2Num) Ast() (functions.Expression, error) {
	op := g.Op2Num[:2]
	num, err := strconv.ParseFloat(g.Op2Num[2:], 64)
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
	Pos    lexer.Position
	Op1Lid string `@Op1Lid`
}

func (g *Op1Lid) Ast() (functions.Expression, error) {
	op := g.Op1Lid[:1]
	name := g.Op1Lid[1:]
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
	Pos    lexer.Position
	Op2Lid string `@Op2Lid`
}

func (g *Op2Lid) Ast() (functions.Expression, error) {
	op := g.Op2Lid[:2]
	name := g.Op2Lid[2:]
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
	NameRegexp string `@NameRegexp`
}

func (g *NameRegexp) Ast() (functions.Expression, error) {
	i := strings.Index(g.NameRegexp, "~")
	name := g.NameRegexp[:i]
	regexpString := g.NameRegexp[i+1 : len(g.NameRegexp)-1]
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
	NameLpar string         `@NameLpar`
	Arg      *Composition   `@@`
	Args     []*Composition `( "," @@ )* ")"`
}

func (g *NameArglist) Ast() (functions.Expression, error) {
	name := g.NameLpar[:len(g.NameLpar)-1]
	args := make([]functions.Expression, len(g.Args)+1)
	var err error
	args[0], err = g.Arg.Ast()
	if err != nil {
		return nil, err
	}
	for i, arg := range g.Args {
		args[i+1], err = arg.Ast()
		if err != nil {
			return nil, err
		}
	}
	return &functions.CallExpression{
		Pos:  g.Pos,
		Name: name,
		Args: args,
	}, nil
}

type NameArray struct {
	Pos        lexer.Position
	NameLbrack string         `@NameLbrack`
	Element    *Composition   `( @@`
	Elements   []*Composition `  ( "," @@ )* )? "]"`
}

func (g *NameArray) Ast() (functions.Expression, error) {
	name := g.NameLbrack[:len(g.NameLbrack)-1]
	arrPos := g.Pos
	arrPos.Column += len(name)
	var elements []functions.Expression
	if g.Element != nil {
		elements = make([]functions.Expression, 1+len(g.Elements))
		var err error
		elements[0], err = g.Element.Ast()
		if err != nil {
			return nil, err
		}
		for i, element := range g.Elements {
			elements[i+1], err = element.Ast()
			if err != nil {
				return nil, err
			}
		}
	}
	return &functions.CallExpression{
		Pos:  g.Pos,
		Name: name,
		Args: []functions.Expression{
			&functions.ArrExpression{
				Pos:      arrPos,
				Elements: elements,
			},
		},
	}, nil
}

type NameObject struct {
	Pos        lexer.Position
	NameLbrace string         `@NameLbarace`
	Prop       *string        `( ( @Lid | @Op1 | @Op2 | @Num )`
	Value      *Composition   `  ":" @@`
	Props      []string       `  ( "," ( @Lid | @Op1 | @Op2 | @Num )`
	Values     []*Composition `    ":" @@ )* )? "}"`
}

func (g *NameObject) Ast() (functions.Expression, error) {
	name := g.NameLbrace[:len(g.NameLbrace)-1]
	objPos := g.Pos
	objPos.Column += len(name)
	propValMap := make(map[string]functions.Expression)
	if g.Prop != nil {
		var err error
		propValMap[*g.Prop], err = g.Value.Ast()
		if err != nil {
			return nil, err
		}
		for i := range g.Props {
			propValMap[g.Props[i]], err = g.Values[i].Ast()
			if err != nil {
				return nil, err
			}
		}
	}
	return &functions.CallExpression{
		Name: name,
		Args: []functions.Expression{
			&functions.ObjExpression{
				objPos,
				propValMap,
			},
		},
	}, nil
}

type NameString struct {
	Pos     lexer.Position
	NameStr string `@NameStr`
}

func (g *NameString) Ast() (functions.Expression, error) {
	i := strings.Index(g.NameStr, "\"")
	name := g.NameStr[:i]
	str, err := strconv.Unquote(g.NameStr[i : len(g.NameStr)-1])
	if err != nil {
		return nil, err
	}
	strPos := g.Pos
	strPos.Column += len(name)
	return &functions.CallExpression{
		Pos:  g.Pos,
		Name: name,
		Args: []functions.Expression{
			&functions.ConstantExpression{
				Pos:   strPos,
				Type:  types.StrType{},
				Value: values.StrValue(str),
			},
		},
	}, nil
}
