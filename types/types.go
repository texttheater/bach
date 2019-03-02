package types

import (
	"bytes"
	"fmt"
	"sort"
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

var BoolType = &boolType{}

type boolType struct {
}

func (t *boolType) Subsumes(other Type) bool {
	_, ok := other.(*boolType)
	return ok
}

func (t *boolType) String() string {
	return "Bool"
}

func (t *boolType) ElementType() Type {
	panic("Bool is not a sequence type")
}

///////////////////////////////////////////////////////////////////////////////

var NumType = &numType{}

type numType struct {
}

func (t *numType) Subsumes(other Type) bool {
	_, ok := other.(*numType)
	return ok
}

func (t *numType) String() string {
	return "Num"
}

func (t *numType) ElementType() Type {
	panic("Num is not a sequence type")
}

///////////////////////////////////////////////////////////////////////////////

var StrType = &strType{}

type strType struct {
}

func (t *strType) Subsumes(other Type) bool {
	_, ok := other.(*strType)
	return ok
}

func (t *strType) String() string {
	return "Str"
}

func (t *strType) ElementType() Type {
	panic("Str is not a sequence type")
}

///////////////////////////////////////////////////////////////////////////////

type SeqType struct {
	ElType Type
}

// AnySeqType represents Seq<Any>, the type of sequences with any type of
// element. It is provided as a package variable for convenience.
var AnySeqType = &SeqType{AnyType}

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

func (t *ArrType) String() string {
	return fmt.Sprintf("Arr<%s>", t.ElType)
}

func (t *ArrType) ElementType() Type {
	return t.ElType
}

///////////////////////////////////////////////////////////////////////////////

func NewObjType(propTypeMap map[string]Type) Type {
	props := make([]string, 0, len(propTypeMap))
	for k := range propTypeMap {
		props = append(props, k)
	}
	sort.Strings(props)
	return &objType{
		props:       props,
		propTypeMap: propTypeMap,
	}
}

type objType struct {
	props       []string
	propTypeMap map[string]Type
}

func (t *objType) Subsumes(other Type) bool {
	otherObjType, ok := other.(*objType)
	if !ok {
		return false
	}
	for k, v1 := range t.propTypeMap {
		v2, ok := otherObjType.propTypeMap[k]
		if !ok {
			return false
		}
		if !v1.Subsumes(v2) {
			return false
		}
	}
	return true
}

func (t *objType) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString("Obj<")
	if len(t.props) > 0 {
		buffer.WriteString(t.props[0])
		buffer.WriteString(": ")
		buffer.WriteString(t.propTypeMap[t.props[0]].String())
		for _, prop := range t.props {
			typ := t.propTypeMap[prop]
			buffer.WriteString(", ")
			buffer.WriteString(prop)
			buffer.WriteString(": ")
			buffer.WriteString(typ.String())
		}
	}
	buffer.WriteString(">")
	return buffer.String()
}

func (t *objType) ElementType() Type {
	panic(fmt.Sprintf("%s is not a sequence type", t))
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

var AnyType = &anyType{}

type anyType struct {
}

func (t *anyType) Subsumes(other Type) bool {
	return true
}

func (t *anyType) String() string {
	return "Any"
}

func (t *anyType) ElementType() Type {
	panic("Any is not a sequence type")
}

///////////////////////////////////////////////////////////////////////////////
