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

func (t *nearrType) Subsumes(u Type) bool {
	if VoidType.Subsumes(u) {
		return true
	}
	switch u := u.(type) {
	case *nearrType:
		if !t.headType.Subsumes(u.headType) {
			return false
		}
		return t.tailType.Subsumes(u.tailType)
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

func (t *nearrType) ElementType() Type {
	return Union(t.headType, t.tailType.ElementType())
}

func (t *nearrType) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString("Tup<")
	buffer.WriteString(t.headType.String())
	tailType := t.tailType
	for {
		if VoidArrType.Subsumes(tailType) {
			buffer.WriteString(">")
			return buffer.String()
		}
		nearrTailType, ok := tailType.(*nearrType)
		if !ok {
			break
		}
		buffer.WriteString(", ")
		buffer.WriteString(nearrTailType.headType.String())
		tailType = nearrTailType.tailType
	}
	buffer.Reset()
	buffer.WriteString("Nearr<")
	buffer.WriteString(t.headType.String())
	buffer.WriteString(", ")
	buffer.WriteString(t.tailType.String())
	buffer.WriteString(">")
	return buffer.String()
}
