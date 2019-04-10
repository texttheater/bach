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
	if otherArrType, ok := other.(arrType); ok {
		return t.elementType.Subsumes(otherArrType.elementType)
	}
	if otherNearrType, ok := other.(*nearrType); ok {
		if !t.elementType.Subsumes(otherNearrType.headType) {
			return false
		}
		return t.Subsumes(otherNearrType.tailType)
	}
	if otherUnionType, ok := other.(unionType); ok {
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
