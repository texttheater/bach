package types

import (
	"bytes"
)

func NewTup(elementTypes []Type) Type {
	return NewNearr(elementTypes, &Arr{Void{}})
}

func NewNearr(elementTypes []Type, restType Type) Type {
	var t Type = restType
	for i := len(elementTypes) - 1; i >= 0; i-- {
		t = &Nearr{
			Head: elementTypes[i],
			Tail: t,
		}
	}
	return t
}

type Nearr struct {
	Head Type
	Tail Type
}

func (t *Nearr) Subsumes(u Type) bool {
	switch u := u.(type) {
	case Void:
		return true
	case *Nearr:
		return t.Head.Subsumes(u.Head) && t.Tail.Subsumes(u.Tail)
	case Union:
		return u.inverseSubsumes(t)
	default:
		return false
	}
}

func (t *Nearr) Bind(u Type, bindings map[string]Type) bool {
	switch u := u.(type) {
	case Void:
		return true
	case *Nearr:
		return t.Head.Bind(u.Head, bindings) && t.Tail.Bind(u.Tail, bindings)
	case Union:
		return u.inverseBind(t, bindings)
	default:
		return false
	}
}

func (t *Nearr) Instantiate(bindings map[string]Type) Type {
	return &Nearr{
		Head: t.Head.Instantiate(bindings),
		Tail: t.Tail.Instantiate(bindings),
	}
}

func (t *Nearr) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case Void:
		return u, t
	case *Nearr:
		headIntersection, _ := t.Head.Partition(u.Head)
		if (Void{}).Subsumes(headIntersection) {
			return Void{}, t
		}
		tailIntersection, _ := t.Tail.Partition(u.Tail)
		if (Void{}).Subsumes(tailIntersection) {
			return Void{}, t
		}
		intersection := &Nearr{
			Head: headIntersection,
			Tail: tailIntersection,
		}
		if intersection.Subsumes(t) {
			return intersection, Void{}
		}
		return intersection, t
	case *Arr:
		headIntersection, _ := t.Head.Partition(u.El)
		if (Void{}).Subsumes(headIntersection) {
			return Void{}, t
		}
		tailIntersection, _ := t.Tail.Partition(u)
		if (Void{}).Subsumes(tailIntersection) {
			return Void{}, t
		}
		intersection := &Nearr{
			Head: headIntersection,
			Tail: tailIntersection,
		}
		if intersection.Subsumes(t) {
			return intersection, Void{}
		}
		return intersection, t
	case Union:
		return u.inversePartition(t)
	case Any:
		return t, Void{}
	default:
		return Void{}, t
	}
}

func (t Nearr) ElementType() Type {
	return NewUnion(t.Head, t.Tail.ElementType())
}

func (t *Nearr) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString("Arr<")
	buffer.WriteString(t.Head.String())
	tail := t.Tail
Loop:
	for {
		switch t := tail.(type) {
		case *Nearr:
			buffer.WriteString(", ")
			buffer.WriteString(t.Head.String())
			tail = t.Tail
		case *Arr:
			if !(Void{}).Subsumes(t.El) {
				buffer.WriteString(", ")
				buffer.WriteString(t.El.String())
				buffer.WriteString("...")
			}
			break Loop
		default:
			panic("non-array tail")
		}
	}
	buffer.WriteString(">")
	return buffer.String()
}
