package types

import (
	"bytes"
)

type TypeVariable struct {
	Name       string
	UpperBound Type
}

func (t TypeVariable) Subsumes(u Type) bool {
	switch u := u.(type) {
	case VoidType:
		return true
	case TypeVariable:
		return t.Name == u.Name
	default:
		return false
	}
}

func (t TypeVariable) Bind(u Type, bindings map[string]Type) bool {
	instType, ok := bindings[t.Name]
	if !ok {
		if t.UpperBound == nil {
			instType = AnyType{}
		} else {
			instType = t.UpperBound
		}
	}
	var newInstType Type
	// pick the more specific type
	if instType.Subsumes(u) {
		newInstType = u
	} else if u.Subsumes(instType) {
		newInstType = instType
	} else {
		return false
	}
	bindings[t.Name] = newInstType
	return true
}

func (t TypeVariable) Instantiate(bindings map[string]Type) Type {
	instType, ok := bindings[t.Name]
	if !ok {
		return t
	}
	return instType
}

func (t TypeVariable) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case VoidType:
		return u, t
	case TypeVariable:
		if t.Name == u.Name {
			return u, VoidType{}
		}
		return VoidType{}, t
	case UnionType:
		return u.inversePartition(t)
	case AnyType:
		return t, VoidType{}
	default:
		return VoidType{}, t
	}
}

func (t TypeVariable) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString("<")
	buffer.WriteString(t.Name)
	if t.UpperBound != nil {
		buffer.WriteString(" ")
		buffer.WriteString(t.UpperBound.String())
	}
	buffer.WriteString(">")
	return buffer.String()
}

func (t TypeVariable) ElementType() Type {
	panic("type variable " + t.Name + " is not a sequence type")
}
