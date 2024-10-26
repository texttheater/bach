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
		result = types.NewUnion(result, t)
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
	return types.Void{}
}

type NullType struct {
	Pos lexer.Position `"Null"`
}

func (g *NullType) Ast() types.Type {
	return types.Null{}
}

type ReaderType struct {
	Pos lexer.Position `"Reader"`
}

func (g *ReaderType) Ast() types.Type {
	return types.Reader{}
}

type BoolType struct {
	Pos lexer.Position `"Bool"`
}

func (g *BoolType) Ast() types.Type {
	return types.Bool{}
}

type NumType struct {
	Pos lexer.Position `"Num"`
}

func (g *NumType) Ast() types.Type {
	return types.Num{}
}

type StrType struct {
	Pos lexer.Position `"Str"`
}

func (g *StrType) Ast() types.Type {
	return types.Str{}
}

type ArrType struct {
	Pos      lexer.Position `"Arr<"`
	Type     *Type          `( @@`
	Types    []*Type        `  ( "," @@ )*`
	Ellipsis *string        `  @Ellipsis? )? ">"`
}

func (g *ArrType) Ast() types.Type {
	if g.Type == nil {
		return types.VoidArr
	}
	if len(g.Types) == 0 && g.Ellipsis != nil {
		return types.NewArr(g.Type.Ast())
	}
	result := &types.Nearr{
		Head: g.Type.Ast(),
	}
	current := result
	length := len(g.Types)
	for i, t := range g.Types {
		if g.Ellipsis != nil && i == length-1 {
			current.Tail = &types.Arr{
				El: t.Ast(),
			}
			return result
		}
		newTail := &types.Nearr{
			Head: t.Ast(),
		}
		current.Tail = newTail
		current = newTail
	}
	current.Tail = types.VoidArr
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
		restTypeAst = types.Any{}
	} else {
		restTypeAst = restType.Ast()
	}
	return types.Obj{
		Props: propTypeMap,
		Rest:  restTypeAst,
	}
}

type AnyType struct {
	Pos lexer.Position `"Any"`
}

func (g *AnyType) Ast() types.Type {
	return types.Any{}
}

type TypeVariable struct {
	Pos       lexer.Position
	LangleLid string `@LangleLid ">"`
}

func (g *TypeVariable) Ast() types.Type {
	return types.NewVar(
		g.LangleLid[1:len(g.LangleLid)],
		types.Any{},
	)
}
