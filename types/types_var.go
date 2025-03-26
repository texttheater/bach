package types

import (
	"bytes"
)

type TypeVar struct {
	Name  string
	Bound Type
}

func NewTypeVar(name string, bound Type) TypeVar {
	return TypeVar{
		Name:  name,
		Bound: bound,
	}
}

func (t TypeVar) Subsumes(u Type) bool {
	switch u := u.(type) {
	case VoidType:
		return true
	case TypeVar:
		return t.Name == u.Name
	default:
		return false
	}
}

func (t TypeVar) Bind(u Type, bindings map[string]Type) bool {
	if !t.Bound.Subsumes(u) {
		return false
	}
	instType, ok := bindings[t.Name]
	if !ok {
		instType = VoidType{}
	}
	newInstType := NewUnionType(instType, u)
	bindings[t.Name] = newInstType
	return true
}

func (t TypeVar) Instantiate(bindings map[string]Type) Type {
	instType, ok := bindings[t.Name]
	if !ok {
		return t
	}
	return instType
}

func (t TypeVar) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case VoidType:
		return u, t
	case TypeVar:
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

func (t TypeVar) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString("<")
	buffer.WriteString(t.Name)
	if t.Bound != nil {
		buffer.WriteString(" ")
		buffer.WriteString(t.Bound.String())
	}
	buffer.WriteString(">")
	return buffer.String()
}

func (t TypeVar) ElementType() Type {
	panic("type variable " + t.Name + " is not a sequence type")
}
