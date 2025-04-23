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
	Pos           lexer.Position  `"Arr<"`
	TypeTemplate  *TypeTemplate   `( @@`
	TypeTemplates []*TypeTemplate `  ( "," @@ )*`
	Ellipsis      *string         `  @Ellipsis? )? ">"`
}

func (g *ArrTypeTemplate) Ast() types.Type {
	if g.TypeTemplate == nil {
		return types.VoidArr
	}
	if len(g.TypeTemplates) == 0 && g.Ellipsis != nil {
		return types.NewArr(g.TypeTemplate.Ast())
	}
	result := &types.Nearr{
		Head: g.TypeTemplate.Ast(),
	}
	current := result
	length := len(g.TypeTemplates)
	for i, t := range g.TypeTemplates {
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

type ObjTypeTemplate struct {
	Pos       lexer.Position            `"Obj<"`
	Prop      *string                   `( ( @Lid | @Op1 | @Op2 | @NumLiteral )`
	ValType   *TypeTemplate             `  ":" @@`
	AfterProp *ObjTypeTemplateAfterProp `  @@`
	RestType  *TypeTemplate             `| @@? ">" )`
}

type ObjTypeTemplateAfterProp struct {
	Pos       lexer.Position            `( ">"`
	Prop      *string                   `| "," ( ( @Lid | @Op1 | @Op2 | @NumLiteral )`
	ValType   *TypeTemplate             `        ":" @@`
	AfterProp *ObjTypeTemplateAfterProp `        @@`
	RestType  *TypeTemplate             `      | @@ ">" ) )`
}

func (g *ObjTypeTemplate) Ast() types.Type {
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

type TypeVariableTemplate struct {
	Pos        lexer.Position
	LangleULid string `@LangleULid`
	UpperBound *Type  `( @@ )? ">"`
}

func (g *TypeVariableTemplate) Ast() types.Type {
	t := types.Var{
		Name: g.LangleULid[1:len(g.LangleULid)],
	}
	if g.UpperBound != nil {
		t.Bound = g.UpperBound.Ast()
	} else {
		t.Bound = types.Any{}
	}
	return t
}
