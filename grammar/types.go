package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/types"
)

type Type struct {
	Pos                lexer.Position
	NonDisjunctiveType *NonDisjunctiveType   `  @@`
	Disjuncts          []*NonDisjunctiveType `{ "|" @@ }`
}

func (g *Type) Ast() types.Type {
	result := g.NonDisjunctiveType.Ast()
	for _, d := range g.Disjuncts {
		result = types.Union(result, d.Ast())
	}
	return result
}

type NonDisjunctiveType struct {
	Pos      lexer.Position
	VoidType *VoidType `  @@`
	NullType *NullType `| @@`
	BoolType *BoolType `| @@`
	NumType  *NumType  `| @@`
	StrType  *StrType  `| @@`
	SeqType  *SeqType  `| @@`
	ArrType  *ArrType  `| @@`
	ObjType  *ObjType  `| @@`
	AnyType  *AnyType  `| @@`
}

func (g *NonDisjunctiveType) Ast() types.Type {
	if g.VoidType != nil {
		return g.VoidType.Ast()
	}
	if g.NullType != nil {
		return g.NullType.Ast()
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
	if g.SeqType != nil {
		return g.SeqType.Ast()
	}
	if g.ArrType != nil {
		return g.ArrType.Ast()
	}
	if g.ObjType != nil {
		return g.ObjType.Ast()
	}
	if g.AnyType != nil {
		return g.AnyType.Ast()
	}
	panic("invalid type")
}

type VoidType struct {
	Pos lexer.Position `"Void"`
}

func (g *VoidType) Ast() types.Type {
	return types.VoidType
}

type NullType struct {
	Pos lexer.Position `"Null"`
}

func (g *NullType) Ast() types.Type {
	return types.NullType
}

type BoolType struct {
	Pos lexer.Position `"Bool"`
}

func (g *BoolType) Ast() types.Type {
	return types.BoolType
}

type NumType struct {
	Pos lexer.Position `"Num"`
}

func (g *NumType) Ast() types.Type {
	return types.NumType
}

type StrType struct {
	Pos lexer.Position `"Str"`
}

func (g *StrType) Ast() types.Type {
	return types.StrType
}

type SeqType struct {
	Pos         lexer.Position `"Seq<"`
	ElementType *Type          `@@ ">"`
}

func (g *SeqType) Ast() types.Type {
	return types.SeqType(g.ElementType.Ast())
}

type ArrType struct {
	Pos         lexer.Position `"Arr<"`
	ElementType *Type          `@@ ">"`
}

func (g *ArrType) Ast() types.Type {
	return types.ArrType(g.ElementType.Ast())
}

type ObjType struct {
	Pos      lexer.Position `"Obj<"`
	Prop     *string        `[ @Prop`
	ValType  *Type          `  ":" @@`
	Props    []string       `  { @Prop`
	ValTypes []*Type        `     ":" @@ } ] ">"`
}

func (g *ObjType) Ast() types.Type {
	propTypeMap := make(map[string]types.Type)
	if g.Prop != nil {
		propTypeMap[*g.Prop] = g.ValType.Ast()
		for i := range g.Props {
			propTypeMap[g.Props[i]] = g.ValTypes[i].Ast()
		}
	}
	return types.ObjType(propTypeMap)
}

type AnyType struct {
	Pos lexer.Position `"Any"`
}

func (g *AnyType) Ast() types.Type {
	return types.AnyType
}
