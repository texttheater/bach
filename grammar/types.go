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
	Pos      lexer.Position `"Tup<"`
	Type     *Type          `( @@`
	Types    []*Type        `  ( "," @@ )*`
	Ellipsis *string        `  @Ellipsis? )? ">"`
}

func (g *TupType) Ast() types.Type {
	if g.Type == nil {
		return types.VoidArrType
	}
	result := &types.NearrType{
		HeadType: g.Type.Ast(),
	}
	current := result
	length := len(g.Types)
	for i, t := range g.Types {
		if g.Ellipsis != nil && i == length-1 {
			current.TailType = &types.ArrType{
				ElType: t.Ast(),
			}
			return result
		}
		newTail := &types.NearrType{
			HeadType: t.Ast(),
		}
		current.TailType = newTail
		current = newTail
	}
	current.TailType = types.VoidArrType
	return result
}

type ObjType struct {
	Pos       lexer.Position    `"Obj<"`
	Prop      *string           `( ( @Lid | @Op1 | @Op2 | @NumLiteral )`
	ValType   *Type             `  ":" @@`
	AfterProp *ObjTypeAfterProp `   @@`
	RestType  *Type             `| @@? ">" )`
}

type ObjTypeAfterProp struct {
	Pos       lexer.Position    `( ">"`
	Prop      *string           `| "," ( ( @Lid | @Op1 | @Op2 | @NumLiteral )`
	ValType   *Type             `        ":" @@`
	AfterProp *ObjTypeAfterProp `        @@`
	RestType  *Type             `      | @@ ">" ) )`
}

func (g *ObjType) Ast() types.Type {
	propTypeMap := make(map[string]types.Type)
	restType := g.RestType
	if g.Prop != nil {
		propTypeMap[*g.Prop] = g.ValType.Ast()
		afterProp := g.AfterProp
		for afterProp != nil {
			if afterProp.Prop != nil {
				propTypeMap[*afterProp.Prop] = afterProp.ValType.Ast()
			}
			restType = afterProp.RestType
			afterProp = afterProp.AfterProp
		}
	}
	var restTypeAst types.Type
	if restType == nil {
		restTypeAst = types.AnyType{}
	} else {
		restTypeAst = restType.Ast()
	}
	return types.ObjType{
		PropTypeMap: propTypeMap,
		RestType:    restTypeAst,
	}
}

type AnyType struct {
	Pos lexer.Position `"Any"`
}

func (g *AnyType) Ast() types.Type {
	return types.AnyType{}
}

type TypeVariable struct {
	Pos       lexer.Position
	LangleLid string `@LangleLid ">"`
}

func (g *TypeVariable) Ast() types.Type {
	return types.TypeVariable{
		Name: g.LangleLid[1:len(g.LangleLid)],
	}
}
