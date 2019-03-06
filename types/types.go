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

var VoidType = voidType{}

type voidType struct {
}

func (t voidType) Subsumes(other Type) bool {
	_, ok := other.(voidType)
	return ok
}

func (t voidType) String() string {
	return "Void"
}

func (t voidType) ElementType() Type {
	panic("Void is not a sequence type")
}

///////////////////////////////////////////////////////////////////////////////

var NullType = nullType{}

type nullType struct {
}

func (t nullType) Subsumes(other Type) bool {
	if VoidType.Subsumes(other) {
		return true
	}
	_, ok := other.(nullType)
	return ok
}

func (t nullType) String() string {
	return "Null"
}

func (t nullType) ElementType() Type {
	panic("Null is not a sequence type")
}

///////////////////////////////////////////////////////////////////////////////

var BoolType = boolType{}

type boolType struct {
}

func (t boolType) Subsumes(other Type) bool {
	if VoidType.Subsumes(other) {
		return true
	}
	_, ok := other.(boolType)
	return ok
}

func (t boolType) String() string {
	return "Bool"
}

func (t boolType) ElementType() Type {
	panic("Bool is not a sequence type")
}

///////////////////////////////////////////////////////////////////////////////

var NumType = numType{}

type numType struct {
}

func (t numType) Subsumes(other Type) bool {
	if VoidType.Subsumes(other) {
		return true
	}
	_, ok := other.(numType)
	return ok
}

func (t numType) String() string {
	return "Num"
}

func (t numType) ElementType() Type {
	panic("Num is not a sequence type")
}

///////////////////////////////////////////////////////////////////////////////

var StrType = strType{}

type strType struct {
}

func (t strType) Subsumes(other Type) bool {
	if VoidType.Subsumes(other) {
		return true
	}
	_, ok := other.(strType)
	return ok
}

func (t strType) String() string {
	return "Str"
}

func (t strType) ElementType() Type {
	panic("Str is not a sequence type")
}

///////////////////////////////////////////////////////////////////////////////

func SeqType(elementType Type) Type {
	return seqType{elementType}
}

type seqType struct {
	elementType Type
}

// AnySeqType represents Seq<Any>, the type of sequences with any type of
// element. It is provided as a package variable for convenience.
var AnySeqType = SeqType(AnyType)

func (t seqType) Subsumes(other Type) bool {
	if VoidType.Subsumes(other) {
		return true
	}
	otherSeqType, ok := other.(seqType)
	if ok {
		return t.elementType.Subsumes(otherSeqType.elementType)
	}
	otherArrType, ok := other.(arrType)
	if ok {
		return t.elementType.Subsumes(otherArrType.elementType)
	}
	return false
}

func (t seqType) String() string {
	return fmt.Sprintf("Seq<%v>", t.elementType)
}

func (t seqType) ElementType() Type {
	return t.elementType
}

///////////////////////////////////////////////////////////////////////////////

func ArrType(elementType Type) Type {
	return arrType{elementType}
}

type arrType struct {
	elementType Type
}

func (t arrType) Subsumes(other Type) bool {
	if VoidType.Subsumes(other) {
		return true
	}
	otherArrType, ok := other.(arrType)
	if !ok {
		return false
	}
	return t.elementType.Subsumes(otherArrType.elementType)
}

func (t arrType) String() string {
	return fmt.Sprintf("Arr<%s>", t.elementType)
}

func (t arrType) ElementType() Type {
	return t.elementType
}

///////////////////////////////////////////////////////////////////////////////

func ObjType(propTypeMap map[string]Type) Type {
	props := make([]string, len(propTypeMap))
	i := 0
	for k := range propTypeMap {
		props[i] = k
		i++
	}
	sort.Strings(props)
	return objType{
		props:       props,
		propTypeMap: propTypeMap,
	}
}

type objType struct {
	props       []string
	propTypeMap map[string]Type
}

func (t objType) Subsumes(other Type) bool {
	if VoidType.Subsumes(other) {
		return true
	}
	otherObjType, ok := other.(objType)
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

func (t objType) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString("Obj<")
	if len(t.props) > 0 {
		buffer.WriteString(t.props[0])
		buffer.WriteString(": ")
		buffer.WriteString(t.propTypeMap[t.props[0]].String())
		for _, prop := range t.props[1:] {
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

func (t objType) ElementType() Type {
	panic(t.String() + " is not a sequence type")
}

///////////////////////////////////////////////////////////////////////////////

func Disjoin(a Type, b Type) Type {
	aDisj, ok := a.(disjunctiveType)
	if ok {
		return aDisj.disjoin(b)
	}
	bDisj, ok := b.(disjunctiveType)
	if ok {
		return bDisj.disjoin(a)
	}
	if a.Subsumes(b) {
		return a
	}
	if b.Subsumes(a) {
		return b
	}
	return disjunctiveType{[]Type{a, b}}
}

type disjunctiveType struct {
	disjuncts []Type
}

func (t disjunctiveType) Subsumes(other Type) bool {
	otherDisj, ok := other.(disjunctiveType)
	if ok {
		return t.subsumesDisj(otherDisj)
	}
	return t.subsumesNonDisj(other)
}

func (t disjunctiveType) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString(fmt.Sprintf("%s", t.disjuncts[0]))
	for _, disjunct := range t.disjuncts[1:] {
		buffer.WriteString("|")
		buffer.WriteString(fmt.Sprintf("%s", disjunct))
	}
	return buffer.String()
}

func (t disjunctiveType) ElementType() Type {
	panic(t.String() + " is not a sequence type")
}

func (t disjunctiveType) subsumesDisj(other disjunctiveType) bool {
	for _, disjunct := range other.disjuncts {
		if !t.subsumesNonDisj(disjunct) {
			return false
		}
	}
	return true
}

func (t disjunctiveType) subsumesNonDisj(other Type) bool {
	for _, disjunct := range t.disjuncts {
		if disjunct.Subsumes(other) {
			return true
		}
	}
	return false
}

func (t disjunctiveType) disjoin(other Type) Type {
	otherDisj, ok := other.(disjunctiveType)
	if ok {
		return t.disjoinDisj(otherDisj)
	}
	return t.disjoinNonDisj(other)
}

func (t disjunctiveType) disjoinDisj(other disjunctiveType) Type {
	result := t
	for _, disjunct := range other.disjuncts {
		result = result.disjoinNonDisj(disjunct)
	}
	return result
}

func (t disjunctiveType) disjoinNonDisj(other Type) disjunctiveType {
	for _, disjunct := range t.disjuncts {
		if disjunct.Subsumes(other) {
			return t
		}
	}
	newDisjuncts := make([]Type, len(t.disjuncts)+1)
	for i, disjunct := range t.disjuncts {
		if !other.Subsumes(disjunct) {
			newDisjuncts[i] = disjunct
		}
	}
	newDisjuncts[len(t.disjuncts)] = other
	return disjunctiveType{newDisjuncts}
}

///////////////////////////////////////////////////////////////////////////////

var AnyType = anyType{}

type anyType struct {
}

func (t anyType) Subsumes(other Type) bool {
	return true
}

func (t anyType) String() string {
	return "Any"
}

func (t anyType) ElementType() Type {
	panic("Any is not a sequence type")
}

///////////////////////////////////////////////////////////////////////////////
