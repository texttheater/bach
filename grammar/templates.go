package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/types"
)

type TypeTemplate struct {
	Pos                        lexer.Position
	NonDisjunctiveTypeTemplate *NonDisjunctiveTypeTemplate   `@@`
	Disjuncts                  []*NonDisjunctiveTypeTemplate `( "|" @@ )*`
}

func (g *TypeTemplate) Ast() types.Type {
	result := g.NonDisjunctiveTypeTemplate.Ast()
	for _, d := range g.Disjuncts {
		t := d.Ast()
		result = types.NewUnion(result, t)
	}
	return result
}

type NonDisjunctiveTypeTemplate struct {
	Pos                  lexer.Position
	VoidType             *VoidType             `  @@`
	NullType             *NullType             `| @@`
	ReaderType           *ReaderType           `| @@`
	BoolType             *BoolType             `| @@`
	NumType              *NumType              `| @@`
	StrType              *StrType              `| @@`
	ArrTypeTemplate      *ArrTypeTemplate      `| @@`
	TupTypeTemplate      *TupTypeTemplate      `| @@`
	ObjTypeTemplate      *ObjTypeTemplate      `| @@`
	AnyType              *AnyType              `| @@`
	TypeVariableTemplate *TypeVariableTemplate `| @@`
}

func (g *NonDisjunctiveTypeTemplate) Ast() types.Type {
	if g.VoidType != nil {
		return g.VoidType.Ast()
	}
	if g.NullType != nil {
		return g.NullType.Ast()
	}
	if g.ReaderType != nil {
		return g.ReaderType.Ast()
	}
	if g.BoolType != nil {
		return g.BoolType.Ast()
	}
	if g.NumType != nil {
		return g.NumType.Ast()
	}
	if g.StrType != nil {
		return g.StrType.Ast()
	}
	if g.ArrTypeTemplate != nil {
		return g.ArrTypeTemplate.Ast()
	}
	if g.TupTypeTemplate != nil {
		return g.TupTypeTemplate.Ast()
	}
	if g.ObjTypeTemplate != nil {
		return g.ObjTypeTemplate.Ast()
	}
	if g.AnyType != nil {
		return g.AnyType.Ast()
	}
	if g.TypeVariableTemplate != nil {
		return g.TypeVariableTemplate.Ast()
	}
	panic("invalid type")
}

type ArrTypeTemplate struct {
	Pos                 lexer.Position `"Arr<"`
	ElementTypeTemplate *TypeTemplate  `@@ ">"`
}

func (g *ArrTypeTemplate) Ast() types.Type {
	elType := g.ElementTypeTemplate.Ast()
	return &types.Arr{elType}
}

type TupTypeTemplate struct {
	Pos           lexer.Position  `"Tup<"`
	TypeTemplate  *TypeTemplate   `( @@`
	TypeTemplates []*TypeTemplate `  ( "," @@ )* )? ">"`
}

func (g *TupTypeTemplate) Ast() types.Type {
	var elementTypes []types.Type
	if g.TypeTemplate != nil {
		elementTypes = make([]types.Type, len(g.TypeTemplates)+1)
		el := g.TypeTemplate.Ast()
		elementTypes[0] = el
		for i, elementTypeTemplate := range g.TypeTemplates {
			el = elementTypeTemplate.Ast()
			elementTypes[i+1] = el
		}
	}
	return types.NewTup(elementTypes)
}

type ObjTypeTemplate struct {
	Pos               lexer.Position `"Obj<"`
	Prop              *string        `( ( @Lid | @Op1 | @Op2 | @NumLiteral )`
	ValTypeTemplate   *Type          `  ":" @@`
	Props             []string       `  ( ( @Lid | @Op1 | @Op2 | @NumLiteral )`
	ValTypeTemplates  []*Type        `     ":" @@ )*`
	RestTypeTemplate1 *Type          `  ( "," @@ )?`
	RestTypeTemplate2 *Type          `| ( @@ )? ) ">"`
}

func (g *ObjTypeTemplate) Ast() types.Type {
	propTypeMap := make(map[string]types.Type)
	if g.Prop != nil {
		propTypeMap[*g.Prop] = g.ValTypeTemplate.Ast()
		for i := range g.Props {
			propTypeMap[g.Props[i]] = g.ValTypeTemplates[i].Ast()
		}
	}
	var restType types.Type
	if g.RestTypeTemplate1 != nil {
		restType = g.RestTypeTemplate1.Ast()
	} else if g.RestTypeTemplate2 != nil {
		restType = g.RestTypeTemplate2.Ast()
	} else {
		restType = types.Any{}
	}
	return types.Obj{
		Props: propTypeMap,
		Rest:  restType,
	}
}

type TypeVariableTemplate struct {
	Pos        lexer.Position
	LangleLid  string `@LangleLid`
	UpperBound *Type  `( @@ )? ">"`
}

func (g *TypeVariableTemplate) Ast() types.Type {
	t := types.Var{
		Name: g.LangleLid[1:len(g.LangleLid)],
	}
	if g.UpperBound != nil {
		t.Bound = g.UpperBound.Ast()
	}
	return t
}
