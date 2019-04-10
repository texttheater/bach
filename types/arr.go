package types

import (
	"fmt"
)

func ArrType(elementType Type) Type {
	return arrType{elementType}
}

var AnyArrType = ArrType(AnyType)

var VoidArrType = ArrType(VoidType)

type arrType struct {
	elementType Type
}

func (t arrType) Subsumes(u Type) bool {
	if VoidType.Subsumes(u) {
		return true
	}
	switch u := u.(type) {
	case arrType:
		return t.elementType.Subsumes(u.elementType)
	case *nearrType:
		if !t.elementType.Subsumes(u.headType) {
			return false
		}
		return t.Subsumes(u.tailType)
	case unionType:
		for _, disjunct := range u {
			if !t.Subsumes(disjunct) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func (t arrType) String() string {
	return fmt.Sprintf("Arr<%s>", t.elementType)
}

func (t arrType) ElementType() Type {
	return t.elementType
}
