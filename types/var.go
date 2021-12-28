package types

import (
	"bytes"
)

type Var struct {
	Name  string
	Bound Type
}

func NewVar(name string, bound Type) Var {
	return Var{
		Name:  name,
		Bound: bound,
	}
}

func (t Var) Subsumes(u Type) bool {
	switch u := u.(type) {
	case Void:
		return true
	case Var:
		return t.Name == u.Name
	default:
		return false
	}
}

func (t Var) Bind(u Type, bindings map[string]Type) bool {
	if !t.Bound.Subsumes(u) {
		return false
	}
	instType, ok := bindings[t.Name]
	if !ok {
		instType = Void{}
	}
	newInstType := NewUnion(instType, u)
	bindings[t.Name] = newInstType
	return true
}

func (t Var) Instantiate(bindings map[string]Type) Type {
	instType, ok := bindings[t.Name]
	if !ok {
		return t
	}
	return instType
}

func (t Var) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case Void:
		return u, t
	case Var:
		if t.Name == u.Name {
			return u, Void{}
		}
		return Void{}, t
	case Union:
		return u.inversePartition(t)
	case Any:
		return t, Void{}
	default:
		return Void{}, t
	}
}

func (t Var) String() string {
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

func (t Var) ElementType() Type {
	panic("type variable " + t.Name + " is not a sequence type")
}
