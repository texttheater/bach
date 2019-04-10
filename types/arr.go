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

func (t arrType) Subsumes(other Type) bool {
	if VoidType.Subsumes(other) {
		return true
	}
	otherArrType, ok := other.(arrType)
	if ok {
		return t.elementType.Subsumes(otherArrType.elementType)
	}
	otherNearrType, ok := other.(*nearrType)
	if ok {
		if !t.elementType.Subsumes(otherNearrType.headType) {
			return false
		}
		return t.Subsumes(otherNearrType.tailType)
	}
	otherUnionType, ok := other.(unionType)
	if ok {
		for _, disjunct := range otherUnionType {
			if !t.Subsumes(disjunct) {
				return false
			}
		}
		return true
	}
	return false
}

func (t arrType) String() string {
	return fmt.Sprintf("Arr<%s>", t.elementType)
}

func (t arrType) ElementType() Type {
	return t.elementType
}
