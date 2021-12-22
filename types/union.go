package types

import (
	"bytes"
	"fmt"
)

func NewUnion(a Type, b Type) Type {
	if aUnion, ok := a.(Union); ok {
		if bUnion, ok := b.(Union); ok {
			for _, bDisjunct := range bUnion {
				aUnion = typeAppend(aUnion, bDisjunct)
			}
			return aUnion
		}
		return typeAppend(aUnion, b)
	}
	if bUnion, ok := b.(Union); ok {
		return typeAppend(bUnion, a)
	}
	if a.Subsumes(b) {
		return a
	}
	if b.Subsumes(a) {
		return b
	}
	return Union([]Type{a, b})
}

// A Union is a slice of types, representing their union. The elements are
// called "disjuncts". The following invariants must be maintained: for every
// Union, 1) no disjunct subsumes another, 2) no disjunct is itself a union
// type.
type Union []Type

func typeAppend(t Union, u Type) Union {
	for i, disjunct := range t {
		// case 1: a disjunct subsumes u already, no change needed
		if disjunct.Subsumes(u) {
			return t
		}
		// case 2: u subsumes a disjunct, need to rebuild slice
		if u.Subsumes(disjunct) {
			newT := make([]Type, i)
			copy(t, newT)
			newT = append(newT, u)
			for j := i + 1; j < len(t); j++ {
				if !u.Subsumes(t[j]) {
					newT = append(newT, t[j])
				}
			}
			return newT
		}
	}
	// case 3: u is completely new, just append it
	return append(t, u)
}

func (t Union) inverseSubsumes(u Type) bool {
	// precondition: u is not a UnionType
	for _, disjunct := range t {
		if !u.Subsumes(disjunct) {
			return false
		}
	}
	return true
}

func (t Union) inverseBind(u Type, bindings map[string]Type) bool {
	// precondition: u is not a UnionType
	for _, disjunct := range t {
		if !u.Bind(disjunct, bindings) {
			return false
		}
	}
	return true
}

func (t Union) Subsumes(u Type) bool {
	switch u := u.(type) {
	case Union:
		// check that for every disjunct of u, at least one disjunct of
		// t subsumes it
	uDisjuncts:
		for _, uDisjunct := range u {
			for _, tDisjunct := range t {
				if tDisjunct.Subsumes(uDisjunct) {
					continue uDisjuncts
				}
			}
			return false
		}
		return true
	default:
		// check that at least one disjunct of t subsumes u
		for _, tDisjunct := range t {
			if tDisjunct.Subsumes(u) {
				return true
			}
		}
		return false
	}
}

func (t Union) Bind(u Type, bindings map[string]Type) bool {
	switch u := u.(type) {
	case Union:
		// check that for every disjunct of u, at least one disjunct of
		// t subsumes it
	uDisjuncts:
		for _, uDisjunct := range u {
			for _, tDisjunct := range t {
				if tDisjunct.Bind(uDisjunct, bindings) {
					continue uDisjuncts
				}
			}
			return false
		}
		return true
	default:
		// check that at least one disjunct of t subsumes u
		for _, tDisjunct := range t {
			if tDisjunct.Bind(u, bindings) {
				return true
			}
		}
		return false
	}
}

func (t Union) Instantiate(bindings map[string]Type) Type {
	var result Type = Void{}
	for _, disjunct := range t {
		result = NewUnion(result, disjunct.Instantiate(bindings))
	}
	return result
}

func (t Union) inversePartition(u Type) (Type, Type) {
	// precondition: u is not a UnionType
	var intersection Type = Void{}
	complement := u
	for _, disjunct := range t {
		var i Type
		i, complement = complement.Partition(disjunct)
		intersection = NewUnion(intersection, i)
	}
	return intersection, complement
}

func (t Union) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case Void:
		return u, t
	case Union:
		var intersectionsUnion Type = Void{}
		var complementsUnion Type = Void{}
		for _, tDisjunct := range t {
			intersection, complement := u.inversePartition(tDisjunct)
			intersectionsUnion = NewUnion(intersectionsUnion, intersection)
			complementsUnion = NewUnion(complementsUnion, complement)
		}
		return intersectionsUnion, complementsUnion
	case Any:
		return t, Void{}
	default:
		var intersectionsUnion Type = Void{}
		var complementsUnion Type = Void{}
		for _, tDisjunct := range t {
			intersection, complement := tDisjunct.Partition(u)
			intersectionsUnion = NewUnion(intersectionsUnion, intersection)
			complementsUnion = NewUnion(complementsUnion, complement)
		}
		return intersectionsUnion, complementsUnion
	}
}

func (t Union) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString(fmt.Sprintf("%s", t[0]))
	for _, disjunct := range t[1:] {
		buffer.WriteString("|")
		buffer.WriteString(fmt.Sprintf("%s", disjunct))
	}
	return buffer.String()
}

func (t Union) ElementType() Type {
	var elType Type = Void{}
	for _, disjunct := range t {
		elType = NewUnion(elType, disjunct.ElementType())
	}
	return elType
}
