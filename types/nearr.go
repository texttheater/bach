package types

import (
	"bytes"
)

func NearrType(headType Type, tailType Type) Type {
	if !AnyArrType.Subsumes(tailType) {
		panic("tail type must be an array type")
	}
	return &nearrType{headType, tailType}
}

type nearrType struct {
	headType Type
	tailType Type
}

func (t nearrType) Subsumes(other Type) bool {
	if VoidType.Subsumes(other) {
		return true
	}
	otherNearrType, ok := other.(nearrType)
	if !ok {
		return false
	}
	if !t.headType.Subsumes(otherNearrType.headType) {
		return false
	}
	return t.tailType.Subsumes(otherNearrType.tailType)
}

func (t nearrType) ElementType() Type {
	return Disjoin(t.headType, t.tailType.ElementType())
}

func (t nearrType) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString("Nearr<")
	buffer.WriteString(t.headType.String())
	buffer.WriteString(", ")
	buffer.WriteString(t.tailType.String())
	buffer.WriteString(")")
	return buffer.String()
}
