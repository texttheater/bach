package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/types"
)

type Type struct {
	Pos                lexer.Position
	NonDisjunctiveType *NonDisjunctiveType   `@@`
	Disjuncts          []*NonDisjunctiveType `( "|" @@ )*`
}

func (g *Type) Ast() types.Type {
	result := g.NonDisjunctiveType.Ast()
	for _, d := range g.Disjuncts {
		t := d.Ast()
		result = types.Union(result, t)
	}
	return result
}

type NonDisjunctiveType struct {
	Pos          lexer.Position
	VoidType     *VoidType     `  @@`
	NullType     *NullType     `| @@`
	ReaderType   *ReaderType   `| @@`
	BoolType     *BoolType     `| @@`
	NumType      *NumType      `| @@`
	StrType      *StrType      `| @@`
	ArrType      *ArrType      `| @@`
	TupType      *TupType      `| @@`
	ObjType      *ObjType      `| @@`
	AnyType      *AnyType      `| @@`
	TypeVariable *TypeVariable `| @@`
}

func (g *NonDisjunctiveType) Ast() types.Type {
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
	if g.ArrType != nil {
		return g.ArrType.Ast()
	}
	if g.TupType != nil {
		return g.TupType.Ast()
	}
	if g.ObjType != nil {
		return g.ObjType.Ast()
	}
	if g.AnyType != nil {
		return g.AnyType.Ast()
	}
	if g.TypeVariable != nil {
		return g.TypeVariable.Ast()
	}
	panic("invalid type")
}

type VoidType struct {
	Pos lexer.Position `"Void"`
}

func (g *VoidType) Ast() types.Type {
	return types.VoidType{}
}

type NullType struct {
	Pos lexer.Position `"Null"`
}

func (g *NullType) Ast() types.Type {
	return types.NullType{}
}

type ReaderType struct {
	Pos lexer.Position `"Reader"`
}

func (g *ReaderType) Ast() types.Type {
	return types.ReaderType{}
}

type BoolType struct {
	Pos lexer.Position `"Bool"`
}

func (g *BoolType) Ast() types.Type {
	return types.BoolType{}
}

type NumType struct {
	Pos lexer.Position `"Num"`
}

func (g *NumType) Ast() types.Type {
	return types.NumType{}
}

type StrType struct {
	Pos lexer.Position `"Str"`
}

func (g *StrType) Ast() types.Type {
	return types.StrType{}
}

type ArrType struct {
	Pos         lexer.Position `"Arr<"`
	ElementType *Type          `@@ ">"`
}

func (g *ArrType) Ast() types.Type {
	elType := g.ElementType.Ast()
	return &types.ArrType{elType}
}

type TupType struct {
	Pos   lexer.Position `"Tup<"`
	Type  *Type          `( @@`
	Types []*Type        `  ( "," @@ )* )? ">"`
}

func (g *TupType) Ast() types.Type {
	var elementTypes []types.Type
	if g.Type != nil {
		elementTypes = make([]types.Type, len(g.Types)+1)
		el := g.Type.Ast()
		elementTypes[0] = el
		for i, elementType := range g.Types {
			el = elementType.Ast()
			elementTypes[i+1] = el
		}
	}
	return types.TupType(elementTypes)
}

type ObjType struct {
	Pos      lexer.Position `"Obj<"`
	Prop     *string        `( ( @Lid | @Op1 | @Op2 | @Num )`
	ValType  *Type          `  ":" @@`
	Props    []string       `  ( ( @Lid | @Op1 | @Op2 | @Num )`
	ValTypes []*Type        `     ":" @@ )* )? ">"`
}

func (g *ObjType) Ast() types.Type {
	propTypeMap := make(map[string]types.Type)
	if g.Prop != nil {
		propTypeMap[*g.Prop] = g.ValType.Ast()
		for i := range g.Props {
			propTypeMap[g.Props[i]] = g.ValTypes[i].Ast()
		}
	}
	return types.NewObjType(propTypeMap)
}

type AnyType struct {
	Pos lexer.Position `"Any"`
}

func (g *AnyType) Ast() types.Type {
	return types.AnyType{}
}

type TypeVariable struct {
	Pos lexer.Position
	Var string `@Typevar`
}

func (g *TypeVariable) Ast() types.Type {
	return types.TypeVariable{
		Name: g.Var[1 : len(g.Var)-1],
	}
}
