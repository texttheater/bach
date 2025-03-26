package types

import (
	"github.com/alecthomas/participle/lexer"
)

type TypeSyntax struct {
	Pos                lexer.Position
	NonDisjunctiveType *NonDisjunctiveType   `@@`
	Disjuncts          []*NonDisjunctiveType `( "|" @@ )*`
}

func (g *TypeSyntax) Ast() Type {
	result := g.NonDisjunctiveType.Ast()
	for _, d := range g.Disjuncts {
		t := d.Ast()
		result = NewUnionType(result, t)
	}
	return result
}

type NonDisjunctiveType struct {
	Pos                lexer.Position
	VoidTypeSyntax     *VoidTypeSyntax     `  @@`
	NullTypeSyntax     *NullTypeSyntax     `| @@`
	ReaderTypeSyntax   *ReaderTypeSyntax   `| @@`
	BoolTypeSyntax     *BoolTypeSyntax     `| @@`
	NumTypeSyntax      *NumTypeSyntax      `| @@`
	StrTypeSyntax      *StrTypeSyntax      `| @@`
	ArrTypeSyntax      *ArrTypeSyntax      `| @@`
	ObjTypeSyntax      *ObjTypeSyntax      `| @@`
	AnyTypeSyntax      *AnyTypeSyntax      `| @@`
	TypeVariableSyntax *TypeVariableSyntax `| @@`
}

func (g *NonDisjunctiveType) Ast() Type {
	if g.VoidTypeSyntax != nil {
		return g.VoidTypeSyntax.Ast()
	}
	if g.NullTypeSyntax != nil {
		return g.NullTypeSyntax.Ast()
	}
	if g.ReaderTypeSyntax != nil {
		return g.ReaderTypeSyntax.Ast()
	}
	if g.BoolTypeSyntax != nil {
		return g.BoolTypeSyntax.Ast()
	}
	if g.NumTypeSyntax != nil {
		return g.NumTypeSyntax.Ast()
	}
	if g.StrTypeSyntax != nil {
		return g.StrTypeSyntax.Ast()
	}
	if g.ArrTypeSyntax != nil {
		return g.ArrTypeSyntax.Ast()
	}
	if g.ObjTypeSyntax != nil {
		return g.ObjTypeSyntax.Ast()
	}
	if g.AnyTypeSyntax != nil {
		return g.AnyTypeSyntax.Ast()
	}
	if g.TypeVariableSyntax != nil {
		return g.TypeVariableSyntax.Ast()
	}
	panic("invalid type")
}

type VoidTypeSyntax struct {
	Pos lexer.Position `"Void"`
}

func (g *VoidTypeSyntax) Ast() Type {
	return VoidType{}
}

type NullTypeSyntax struct {
	Pos lexer.Position `"Null"`
}

func (g *NullTypeSyntax) Ast() Type {
	return NullType{}
}

type ReaderTypeSyntax struct {
	Pos lexer.Position `"Reader"`
}

func (g *ReaderTypeSyntax) Ast() Type {
	return ReaderType{}
}

type BoolTypeSyntax struct {
	Pos lexer.Position `"Bool"`
}

func (g *BoolTypeSyntax) Ast() Type {
	return BoolType{}
}

type NumTypeSyntax struct {
	Pos lexer.Position `"Num"`
}

func (g *NumTypeSyntax) Ast() Type {
	return NumType{}
}

type StrTypeSyntax struct {
	Pos lexer.Position `"Str"`
}

func (g *StrTypeSyntax) Ast() Type {
	return StrType{}
}

type ArrTypeSyntax struct {
	Pos      lexer.Position `"Arr<"`
	Type     *TypeSyntax    `( @@`
	Types    []*TypeSyntax  `  ( "," @@ )*`
	Ellipsis *string        `  @Ellipsis? )? ">"`
}

func (g *ArrTypeSyntax) Ast() Type {
	if g.Type == nil {
		return VoidArrType
	}
	if len(g.Types) == 0 && g.Ellipsis != nil {
		return NewArrType(g.Type.Ast())
	}
	result := &NearrType{
		Head: g.Type.Ast(),
	}
	current := result
	length := len(g.Types)
	for i, t := range g.Types {
		if g.Ellipsis != nil && i == length-1 {
			current.Tail = &ArrType{
				El: t.Ast(),
			}
			return result
		}
		newTail := &NearrType{
			Head: t.Ast(),
		}
		current.Tail = newTail
		current = newTail
	}
	current.Tail = VoidArrType
	return result
}

type ObjTypeSyntax struct {
	Pos       lexer.Position          `"Obj<"`
	Prop      *string                 `( ( @Lid | @Op1 | @Op2 | @NumLiteral )`
	ValType   *TypeSyntax             `  ":" @@`
	AfterProp *ObjTypeSyntaxAfterProp `   @@`
	RestType  *TypeSyntax             `| @@? ">" )`
}

type ObjTypeSyntaxAfterProp struct {
	Pos       lexer.Position          `( ">"`
	Prop      *string                 `| "," ( ( @Lid | @Op1 | @Op2 | @NumLiteral )`
	ValType   *TypeSyntax             `        ":" @@`
	AfterProp *ObjTypeSyntaxAfterProp `        @@`
	RestType  *TypeSyntax             `      | @@ ">" ) )`
}

func (g *ObjTypeSyntax) Ast() Type {
	propTypeMap := make(map[string]Type)
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
	var restTypeAst Type
	if restType == nil {
		restTypeAst = AnyType{}
	} else {
		restTypeAst = restType.Ast()
	}
	return ObjType{
		Props: propTypeMap,
		Rest:  restTypeAst,
	}
}

type AnyTypeSyntax struct {
	Pos lexer.Position `"Any"`
}

func (g *AnyTypeSyntax) Ast() Type {
	return AnyType{}
}

type TypeVariableSyntax struct {
	Pos       lexer.Position
	LangleLid string `@LangleLid ">"`
}

func (g *TypeVariableSyntax) Ast() Type {
	return NewTypeVar(
		g.LangleLid[1:len(g.LangleLid)],
		AnyType{},
	)
}
