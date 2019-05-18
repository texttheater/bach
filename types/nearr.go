package types

import (
	"bytes"
)

type NearrType struct {
	HeadType Type
	TailType Type
}

func (t *NearrType) Subsumes(u Type) bool {
	switch u := u.(type) {
	case VoidType:
		return true
	case *NearrType:
		if !t.HeadType.Subsumes(u.HeadType) {
			return false
		}
		return t.TailType.Subsumes(u.TailType)
	case UnionType:
		return u.inverseSubsumes(t)
	default:
		return false
	}
}

func (t *NearrType) ElementType() Type {
	return Union(t.HeadType, t.TailType.ElementType())
}

func (t *NearrType) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString("Tup<")
	buffer.WriteString(t.HeadType.String())
	tailType := t.TailType
	for {
		if VoidArrType.Subsumes(tailType) {
			buffer.WriteString(">")
			return buffer.String()
		}
		nearrTailType, ok := tailType.(*NearrType)
		if !ok {
			break
		}
		buffer.WriteString(", ")
		buffer.WriteString(nearrTailType.HeadType.String())
		tailType = nearrTailType.TailType
	}
	buffer.Reset()
	buffer.WriteString("Nearr<")
	buffer.WriteString(t.HeadType.String())
	buffer.WriteString(", ")
	buffer.WriteString(t.TailType.String())
	buffer.WriteString(">")
	return buffer.String()
}
