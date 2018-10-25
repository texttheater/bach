package types

import (
	"fmt"
)

type Type interface {
	Subsumes(Type) bool
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
	return fmt.Sprintf("Seq<%v>", t.ElementType)
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

type ArrayType struct {
	ElementType Type
	HasLength   bool
	Length      uint
}

func (t *ArrayType) Subsumes(other Type) bool {
	otherArrayType, ok := other.(*ArrayType)
	if !ok {
		return false
	}
	if t.HasLength {
		if !otherArrayType.HasLength {
			return false
		}
		if otherArrayType.Length != t.Length {
			return false
		}
	}
	return t.ElementType.Subsumes(otherArrayType.ElementType)
}

func (t ArrayType) String() string {
	return "Arr" // TODO be more specific
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
