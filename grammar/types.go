package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/types"
)

type Type struct {
	Pos      lexer.Position
	NullType *NullType `  @@`
	BoolType *BoolType `| @@`
	NumType  *NumType  `| @@`
	StrType  *StrType  `| @@`
	SeqType  *SeqType  `| @@`
	ArrType  *ArrType  `| @@`
	ObjType  *ObjType  `| @@`
	//DisjunctiveType *DisjunctiveType `| @@` // FIXME parser doesn't handle them
	AnyType *AnyType `| @@`
}

func (g *Type) Ast() types.Type {
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
	//if g.DisjunctiveType != nil {
	//	return g.DisjunctiveType.Ast()
	//}
	if g.AnyType != nil {
		return g.AnyType.Ast()
	}
	panic("invalid type")
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
	Prop     *string        `[ @Name`
	ValType  *Type          `  ":" @@`
	Props    []string       `  { @Name`
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

type DisjunctiveType struct {
	Pos   lexer.Position
	Type1 *Type   `    @@`
	Type2 *Type   `"|" @@`
	Types []*Type `{ "|" @@ }`
}

func (g *DisjunctiveType) Ast() types.Type {
	result := types.Disjoin(g.Type1.Ast(), g.Type2.Ast())
	for _, t := range g.Types {
		result = types.Disjoin(result, t.Ast())
	}
	return result
}

type AnyType struct {
	Pos lexer.Position `"Any"`
}

func (g *AnyType) Ast() types.Type {
	return types.AnyType
}
