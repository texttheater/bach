package types

import (
	"fmt"
)

type Arr struct {
	El Type
}

var AnyArr Type = &Arr{Any{}}

var VoidArr Type = &Arr{Void{}}

func NewArr(el Type) *Arr {
	return &Arr{
		El: el,
	}
}

func (t *Arr) Subsumes(u Type) bool {
	switch u := u.(type) {
	case Void:
		return true
	case *Nearr:
		return t.El.Subsumes(u.Head) && t.Subsumes(u.Tail)
	case *Arr:
		return t.El.Subsumes(u.El)
	case Union:
		return u.inverseSubsumes(t)
	default:
		return false
	}
}

func (t *Arr) Bind(u Type, bindings map[string]Type) bool {
	switch u := u.(type) {
	case Void:
		return true
	case *Nearr:
		return t.El.Bind(u.ElementType(), bindings)
	case *Arr:
		return t.El.Bind(u.El, bindings)
	case Union:
		return u.inverseBind(t, bindings)
	default:
		return false
	}
}

func (t *Arr) Instantiate(bindings map[string]Type) Type {
	return &Arr{t.El.Instantiate(bindings)}
}

func (t *Arr) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case Void:
		return u, t
	case *Nearr:
		headIntersection, _ := t.El.Partition(u.Head)
		if (Void{}).Subsumes(headIntersection) {
			return Void{}, t
		}
		tailIntersection, _ := t.Partition(u.Tail)
		if (Void{}).Subsumes(tailIntersection) {
			return Void{}, t
		}
		return &Nearr{
			Head: headIntersection,
			Tail: tailIntersection,
		}, t
	case *Arr:
		intersection, _ := t.El.Partition(u.El)
		if intersection.Subsumes(t.El) {
			return &Arr{intersection}, Void{}
		}
		return &Arr{intersection}, t
	case Union:
		return u.inversePartition(t)
	case Any:
		return t, Void{}
	default:
		return Void{}, t
	}
}

func (t *Arr) String() string {
	if (Void{}).Subsumes(t.El) {
		return "Tup<>"
	}
	return fmt.Sprintf("Arr<%s>", t.El)
}

func (t *Arr) ElementType() Type {
	return t.El
}
