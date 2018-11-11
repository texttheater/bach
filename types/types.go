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

type BooleanType struct {
}

func (t *BooleanType) Subsumes(other Type) bool {
	_, ok := other.(*BooleanType)
	return ok
}

func (t *BooleanType) String() string {
	return "Bool"
}

///////////////////////////////////////////////////////////////////////////////

type NumberType struct {
}

func (t *NumberType) Subsumes(other Type) bool {
	_, ok := other.(*NumberType)
	return ok
}

func (t *NumberType) String() string {
	return "Num"
}

///////////////////////////////////////////////////////////////////////////////

type StringType struct {
}

func (t *StringType) Subsumes(other Type) bool {
	_, ok := other.(*StringType)
	return ok
}

func (t *StringType) String() string {
	return "Str"
}

///////////////////////////////////////////////////////////////////////////////

type SeqType struct {
	ElementType Type
}

func (t *SeqType) Subsumes(other Type) bool {
	otherSeqType, ok := other.(*SeqType)
	if !ok {
		return false
	}
	return t.ElementType.Subsumes(otherSeqType.ElementType)
}

func (t *SeqType) String() string {
	return fmt.Sprintf("Seq<%s>", t.ElementType)
}

///////////////////////////////////////////////////////////////////////////////

type ArrayType struct {
	ElementType Type
}

func (t *ArrayType) Subsumes(other Type) bool {
	otherArrayType, ok := other.(*ArrayType)
	if !ok {
		return false
	}
	return t.ElementType.Subsumes(otherArrayType.ElementType)
}

func (t ArrayType) String() string {
	return fmt.Sprintf("Arr<%s>", t.ElementType)
}

///////////////////////////////////////////////////////////////////////////////

type ObjectType struct {
	fieldTypeMap map[string]Type
}

func (t *ObjectType) Subsumes(other Type) bool {
	otherObjectType, ok := other.(*ObjectType)
	if !ok {
		return false
	}
	for field, b := range otherObjectType.fieldTypeMap {
		if a, ok := t.fieldTypeMap[field]; ok {
			if !a.Subsumes(b) {
				return false
			}
		}
	}
	return true
}

func (t *ObjectType) String() string {
	return "Obj" // TODO be more specific
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
