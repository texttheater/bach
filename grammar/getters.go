package grammar

import (
	"strconv"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
)

type Getter struct {
	Pos       lexer.Position
	LidGetter *string `  @LidGetter`
	Op1Getter *string `| @Op1Getter`
	Op2Getter *string `| @Op2Getter`
	NumGetter *string `| @NumGetter`
	StrGetter *string `| @StrGetter`
}

func (g *Getter) Ast() (expressions.Expression, error) {
	var name string
	var err error
	if g.LidGetter != nil {
		name = (*g.LidGetter)[1:]
	} else if g.Op1Getter != nil {
		name = (*g.Op1Getter)[1:]
	} else if g.Op2Getter != nil {
		name = (*g.Op2Getter)[1:]
	} else if g.NumGetter != nil {
		numString := (*g.NumGetter)[1:]
		sign := 1.0
		if numString[0] == '-' {
			numString = numString[1:]
			sign = -1.0
		}
		num, err := strconv.ParseFloat(numString, 64)
		if err != nil {
			return nil, err
		}
		num = num * sign
		name = strconv.FormatFloat(float64(num), 'g', -1, 64)
	} else if g.StrGetter != nil {
		name, err = strconv.Unquote((*g.StrGetter)[1:])
		if err != nil {
			return nil, err
		}
	} else {
		panic("invalid getter")
	}
	return &expressions.GetterExpression{
		Pos:  g.Pos,
		Name: name,
	}, nil
}
