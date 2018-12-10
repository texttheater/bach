package types

import (
	"bytes"
	"fmt"
)

///////////////////////////////////////////////////////////////////////////////

type Type interface {
	Subsumes(Type) bool
}

///////////////////////////////////////////////////////////////////////////////

type NullType struct {
}

func (t *NullType) Subsumes(other Type) bool {
	_, ok := other.(*NullType)
	return ok
}

func (t *NullType) String() string {
	return "Null"
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

///////////////////////////////////////////////////////////////////////////////

type SeqType struct {
	ElementType Type
}

func (t *SeqType) Subsumes(other Type) bool {
	otherSeqType, ok := other.(*SeqType)
	if ok {
		return t.ElementType.Subsumes(otherSeqType.ElementType)
	}
	otherArrType, ok := other.(*ArrType)
	if ok {
		return t.ElementType.Subsumes(otherArrType.ElementType)
	}
	return false
}

func (t *SeqType) String() string {
	return fmt.Sprintf("Seq<%s>", t.ElementType)
}

///////////////////////////////////////////////////////////////////////////////

type ArrType struct {
	ElementType Type
}

func (t *ArrType) Subsumes(other Type) bool {
	otherArrType, ok := other.(*ArrType)
	if !ok {
		return false
	}
	return t.ElementType.Subsumes(otherArrType.ElementType)
}

func (t ArrType) String() string {
	return fmt.Sprintf("Arr<%s>", t.ElementType)
}

///////////////////////////////////////////////////////////////////////////////

type DisjunctiveType struct {
	disjuncts []Type
}

func (t *DisjunctiveType) Subsumes(other Type) bool {
	otherDisj, ok := other.(*DisjunctiveType)
	if ok {
		return t.subsumesDisj(otherDisj)
	}
	return t.subsumesNonDisj(other)
}

func (t *DisjunctiveType) subsumesDisj(other *DisjunctiveType) bool {
	for _, disjunct := range other.disjuncts {
		if !t.subsumesNonDisj(disjunct) {
			return false
		}
	}
	return true
}

func (t *DisjunctiveType) subsumesNonDisj(other Type) bool {
	for _, disjunct := range t.disjuncts {
		if disjunct.Subsumes(other) {
			return true
		}
	}
	return false
}

func (t *DisjunctiveType) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString(fmt.Sprintf("%s", t.disjuncts[0]))
	for _, disjunct := range t.disjuncts[1:] {
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
	for _, disjunct := range other.disjuncts {
		result = result.disjoinNonDisj(disjunct)
	}
	return result
}

func (t *DisjunctiveType) disjoinNonDisj(other Type) *DisjunctiveType {
	for _, disjunct := range t.disjuncts {
		if disjunct.Subsumes(other) {
			return t
		}
	}
	newDisjuncts := make([]Type, 0, len(t.disjuncts)+1)
	for _, disjunct := range t.disjuncts {
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

///////////////////////////////////////////////////////////////////////////////
