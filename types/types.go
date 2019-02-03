package types

import (
	"bytes"
	"fmt"
)

///////////////////////////////////////////////////////////////////////////////

type Type interface {
	Subsumes(Type) bool
	String() string
	ElementType() Type
}

///////////////////////////////////////////////////////////////////////////////

var NullType = &nullType{}

type nullType struct {
}

func (t *nullType) Subsumes(other Type) bool {
	_, ok := other.(*nullType)
	return ok
}

func (t *nullType) String() string {
	return "Null"
}

func (t *nullType) ElementType() Type {
	panic("Null is not a sequence type")
}

///////////////////////////////////////////////////////////////////////////////

type BoolType struct {
}

func (t *BoolType) Subsumes(other Type) bool {
	_, ok := other.(*BoolType)
	return ok
}

func (t *BoolType) String() string {
	return "Bool"
}

func (t *BoolType) ElementType() Type {
	panic("Bool is not a sequence type")
}

///////////////////////////////////////////////////////////////////////////////

type NumType struct {
}

func (t *NumType) Subsumes(other Type) bool {
	_, ok := other.(*NumType)
	return ok
}

func (t *NumType) String() string {
	return "Num"
}

func (t *NumType) ElementType() Type {
	panic("Num is not a sequence type")
}

///////////////////////////////////////////////////////////////////////////////

type StrType struct {
}

func (t *StrType) Subsumes(other Type) bool {
	_, ok := other.(*StrType)
	return ok
}

func (t *StrType) String() string {
	return "Str"
}

func (t *StrType) ElementType() Type {
	panic("Str is not a sequence type")
}

///////////////////////////////////////////////////////////////////////////////

type SeqType struct {
	ElType Type
}

func (t *SeqType) Subsumes(other Type) bool {
	otherSeqType, ok := other.(*SeqType)
	if ok {
		return t.ElType.Subsumes(otherSeqType.ElType)
	}
	otherArrType, ok := other.(*ArrType)
	if ok {
		return t.ElType.Subsumes(otherArrType.ElType)
	}
	return false
}

func (t *SeqType) String() string {
	return fmt.Sprintf("Seq<%v>", t.ElType)
}

func (t *SeqType) ElementType() Type {
	return t.ElType
}

///////////////////////////////////////////////////////////////////////////////

type ArrType struct {
	ElType Type
}

func (t *ArrType) Subsumes(other Type) bool {
	otherArrType, ok := other.(*ArrType)
	if !ok {
		return false
	}
	return t.ElType.Subsumes(otherArrType.ElType)
}

func (t ArrType) String() string {
	return fmt.Sprintf("Arr<%s>", t.ElType)
}

func (t *ArrType) ElementType() Type {
	return t.ElType
}

///////////////////////////////////////////////////////////////////////////////

type DisjunctiveType struct {
	Disjuncts []Type
}

func (t *DisjunctiveType) Subsumes(other Type) bool {
	otherDisj, ok := other.(*DisjunctiveType)
	if ok {
		return t.subsumesDisj(otherDisj)
	}
	return t.subsumesNonDisj(other)
}

func (t *DisjunctiveType) ElementType() Type {
	panic(t.String() + " is not a sequence type")
}

func (t *DisjunctiveType) subsumesDisj(other *DisjunctiveType) bool {
	for _, disjunct := range other.Disjuncts {
		if !t.subsumesNonDisj(disjunct) {
			return false
		}
	}
	return true
}

func (t *DisjunctiveType) subsumesNonDisj(other Type) bool {
	for _, disjunct := range t.Disjuncts {
		if disjunct.Subsumes(other) {
			return true
		}
	}
	return false
}

func (t *DisjunctiveType) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString(fmt.Sprintf("%s", t.Disjuncts[0]))
	for _, disjunct := range t.Disjuncts[1:] {
		buffer.WriteString("|")
		buffer.WriteString(fmt.Sprintf("%s", disjunct))
	}
	return buffer.String()
}

func Disjoin(a Type, b Type) Type {
	aDisj, ok := a.(*DisjunctiveType)
	if ok {
		return aDisj.disjoin(b)
	}
	bDisj, ok := b.(*DisjunctiveType)
	if ok {
		return bDisj.disjoin(a)
	}
	if a.Subsumes(b) {
		return a
	}
	if b.Subsumes(a) {
		return b
	}
	return &DisjunctiveType{[]Type{a, b}}
}

func (t *DisjunctiveType) disjoin(other Type) Type {
	otherDisj, ok := other.(*DisjunctiveType)
	if ok {
		return t.disjoinDisj(otherDisj)
	}
	return t.disjoinNonDisj(other)
}

func (t *DisjunctiveType) disjoinDisj(other *DisjunctiveType) Type {
	result := t
	for _, disjunct := range other.Disjuncts {
		result = result.disjoinNonDisj(disjunct)
	}
	return result
}

func (t *DisjunctiveType) disjoinNonDisj(other Type) *DisjunctiveType {
	for _, disjunct := range t.Disjuncts {
		if disjunct.Subsumes(other) {
			return t
		}
	}
	newDisjuncts := make([]Type, 0, len(t.Disjuncts)+1)
	for _, disjunct := range t.Disjuncts {
		if !other.Subsumes(disjunct) {
			newDisjuncts = append(newDisjuncts, disjunct)
		}
	}
	newDisjuncts = append(newDisjuncts, other)
	return &DisjunctiveType{newDisjuncts}
}

///////////////////////////////////////////////////////////////////////////////

type AnyType struct {
}

func (t *AnyType) Subsumes(other Type) bool {
	return true
}

func (t *AnyType) String() string {
	return "Any"
}

func (t *AnyType) ElementType() Type {
	panic("Any is not a sequence type")
}

///////////////////////////////////////////////////////////////////////////////
