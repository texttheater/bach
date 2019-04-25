package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/types"
)

type Type struct {
	Pos                lexer.Position
	NonDisjunctiveType *NonDisjunctiveType   `@@`
	Disjuncts          []*NonDisjunctiveType `( "|" @@ )*`
}

func (g *Type) Ast() (types.Type, error) {
	result, err := g.NonDisjunctiveType.Ast()
	if err != nil {
		return nil, err
	}
	for _, d := range g.Disjuncts {
		t, err := d.Ast()
		if err != nil {
			return nil, err
		}
		result = types.Union(result, t)
	}
	return result, nil
}

type NonDisjunctiveType struct {
	Pos        lexer.Position
	VoidType   *VoidType   `  @@`
	NullType   *NullType   `| @@`
	ReaderType *ReaderType `| @@`
	BoolType   *BoolType   `| @@`
	NumType    *NumType    `| @@`
	StrType    *StrType    `| @@`
	SeqType    *SeqType    `| @@`
	ArrType    *ArrType    `| @@`
	NearrType  *NearrType  `| @@`
	TupType    *TupType    `| @@`
	ObjType    *ObjType    `| @@`
	AnyType    *AnyType    `| @@`
}

func (g *NonDisjunctiveType) Ast() (types.Type, error) {
	if g.VoidType != nil {
		return g.VoidType.Ast(), nil
	}
	if g.NullType != nil {
		return g.NullType.Ast(), nil
	}
	if g.ReaderType != nil {
		return g.ReaderType.Ast(), nil
	}
	if g.BoolType != nil {
		return g.BoolType.Ast(), nil
	}
	if g.NumType != nil {
		return g.NumType.Ast(), nil
	}
	if g.StrType != nil {
		return g.StrType.Ast(), nil
	}
	if g.SeqType != nil {
		return g.SeqType.Ast()
	}
	if g.ArrType != nil {
		return g.ArrType.Ast()
	}
	if g.NearrType != nil {
		ast, err := g.NearrType.Ast()
		if err != nil {
			return nil, err
		}
		return ast, nil
	}
	if g.TupType != nil {
		return g.TupType.Ast()
	}
	if g.ObjType != nil {
		return g.ObjType.Ast()
	}
	if g.AnyType != nil {
		return g.AnyType.Ast(), nil
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

type SeqType struct {
	Pos         lexer.Position `"Seq<"`
	ElementType *Type          `@@ ">"`
}

func (g *SeqType) Ast() (types.Type, error) {
	elType, err := g.ElementType.Ast()
	if err != nil {
		return nil, err
	}
	return &types.SeqType{elType}, nil
}

type ArrType struct {
	Pos         lexer.Position `"Arr<"`
	ElementType *Type          `@@ ">"`
}

func (g *ArrType) Ast() (types.Type, error) {
	elType, err := g.ElementType.Ast()
	if err != nil {
		return nil, err
	}
	return &types.ArrType{elType}, nil
}

type NearrType struct {
	Pos      lexer.Position `"Nearr<"`
	HeadType *Type          `@@ ","`
	TailType *Type          `@@ ">"`
}

func (g *NearrType) Ast() (types.Type, error) {
	headType, err := g.HeadType.Ast()
	if err != nil {
		return nil, err
	}
	tailType, err := g.TailType.Ast()
	if err != nil {
		return nil, err
	}
	if !types.AnyArrType.Subsumes(tailType) {
		return nil, errors.E(
			errors.Code(errors.TailRequiresArrType),
			errors.Pos(g.Pos),
			errors.WantType(types.AnyArrType),
			errors.GotType(tailType),
		)
	}
	return &types.NearrType{headType, tailType}, nil
}

type TupType struct {
	Pos   lexer.Position `"Tup<"`
	Type  *Type          `( @@`
	Types []*Type        `  ( "," @@ )* )? ">"`
}

func (g *TupType) Ast() (types.Type, error) {
	t := types.VoidArrType
	for i := len(g.Types) - 1; i >= 0; i-- {
		headType, err := g.Types[i].Ast()
		if err != nil {
			return nil, err
		}
		t = &types.NearrType{headType, t}
	}
	if g.Type != nil {
		headType, err := g.Type.Ast()
		if err != nil {
			return nil, err
		}
		t = &types.NearrType{headType, t}
	}
	return t, nil
}

type ObjType struct {
	Pos      lexer.Position `"Obj<"`
	Prop     *string        `( @Prop`
	ValType  *Type          `  ":" @@`
	Props    []string       `  ( @Prop`
	ValTypes []*Type        `     ":" @@ )* )? ">"`
}

func (g *ObjType) Ast() (types.Type, error) {
	propTypeMap := make(map[string]types.Type)
	if g.Prop != nil {
		var err error
		propTypeMap[*g.Prop], err = g.ValType.Ast()
		if err != nil {
			return nil, err
		}
		for i := range g.Props {
			propTypeMap[g.Props[i]], err = g.ValTypes[i].Ast()
			if err != nil {
				return nil, err
			}
		}
	}
	return types.NewObjType(propTypeMap), nil
}

type AnyType struct {
	Pos lexer.Position `"Any"`
}

func (g *AnyType) Ast() types.Type {
	return types.AnyType{}
}
