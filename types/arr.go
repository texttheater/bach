package types

import (
	"fmt"
)

type ArrType struct {
	ElType Type
}

var AnyArrType Type = &ArrType{AnyType{}}

var VoidArrType Type = &ArrType{VoidType{}}

func (t *ArrType) Subsumes(u Type) bool {
	switch u := u.(type) {
	case VoidType:
		return true
	case TupType:
		for _, el := range u {
			if !t.ElType.Subsumes(el) {
				return false
			}
		}
		return true
	case *ArrType:
		return t.ElType.Subsumes(u.ElType)
	case UnionType:
		return u.inverseSubsumes(t)
	default:
		return false
	}
}

func (t *ArrType) String() string {
	if (VoidType{}).Subsumes(t.ElType) {
		return "Tup<>"
	}
	return fmt.Sprintf("Arr<%s>", t.ElType)
}

func (t *ArrType) ElementType() Type {
	return t.ElType
}
