package grammar

import (
	"strconv"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
)

type Getter struct {
	Pos       lexer.Position
	LidGetter *string `  @LidGetter`
	Op1Getter *string `| @Op1Getter`
	Op2Getter *string `| @Op2Getter`
	NumGetter *string `| @NumGetter`
	StrGetter *string `| @StrGetter`
}

func (g *Getter) Ast() (functions.Expression, error) {
	var name string
	var err error
	if g.LidGetter != nil {
		name = (*g.LidGetter)[1:]
	} else if g.Op1Getter != nil {
		name = (*g.Op1Getter)[1:]
	} else if g.Op2Getter != nil {
		name = (*g.Op2Getter)[1:]
	} else if g.NumGetter != nil {
		num, err := strconv.ParseFloat((*g.NumGetter)[1:], 64)
		if err != nil {
			return nil, err
		}
		name = strconv.FormatFloat(float64(num), 'g', -1, 64)
	} else if g.StrGetter != nil {
		name, err = strconv.Unquote((*g.StrGetter)[1:])
		if err != nil {
			return nil, err
		}
	} else {
		panic("invalid getter")
	}
	return &functions.GetterExpression{
		Pos:  g.Pos,
		Name: name,
	}, nil
}
