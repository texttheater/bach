package types

import (
	"bytes"
	"fmt"
)

func Union(a Type, b Type) Type {
	if aUnion, ok := a.(UnionType); ok {
		if bUnion, ok := b.(UnionType); ok {
			for _, bDisjunct := range bUnion {
				aUnion = typeAppend(aUnion, bDisjunct)
			}
			return aUnion
		}
		return typeAppend(aUnion, b)
	}
	if bUnion, ok := b.(UnionType); ok {
		return typeAppend(bUnion, a)
	}
	if a.Subsumes(b) {
		return a
	}
	if b.Subsumes(a) {
		return b
	}
	return UnionType([]Type{a, b})
}

// A UnionType is a slice of types, representing their union. The elements are
// called "disjuncts". The following invariants must be maintained: for every
// UnionType, 1) no disjunct subsumes another, 2) no disjunct is itself a union
// type.
type UnionType []Type

func typeAppend(t UnionType, u Type) UnionType {
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

func (t UnionType) inverseSubsumes(u Type) bool {
	// precondition: u is not a UnionType
	for _, disjunct := range t {
		if !u.Subsumes(disjunct) {
			return false
		}
	}
	return true
}

func (t UnionType) Subsumes(u Type) bool {
	switch u := u.(type) {
	case UnionType:
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

func (t UnionType) inversePartition(u Type) (Type, Type) {
	// precondition: u is not a UnionType
	var intersection Type = VoidType{}
	complement := u
	for _, disjunct := range t {
		var i Type
		i, complement = complement.Partition(disjunct)
		intersection = Union(intersection, i)
	}
	return intersection, complement
}

func (t UnionType) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case VoidType:
		return u, t
	case UnionType:
		var intersectionsUnion Type = VoidType{}
		var complementsUnion Type = VoidType{}
		for _, tDisjunct := range t {
			intersection, complement := u.inversePartition(tDisjunct)
			intersectionsUnion = Union(intersectionsUnion, intersection)
			complementsUnion = Union(complementsUnion, complement)
		}
		return intersectionsUnion, complementsUnion
	case AnyType:
		return t, VoidType{}
	default:
		var intersectionsUnion Type = VoidType{}
		var complementsUnion Type = VoidType{}
		for _, tDisjunct := range t {
			intersection, complement := tDisjunct.Partition(u)
			intersectionsUnion = Union(intersectionsUnion, intersection)
			complementsUnion = Union(complementsUnion, complement)
		}
		return intersectionsUnion, complementsUnion
	}
}

func (t UnionType) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString(fmt.Sprintf("%s", t[0]))
	for _, disjunct := range t[1:] {
		buffer.WriteString("|")
		buffer.WriteString(fmt.Sprintf("%s", disjunct))
	}
	return buffer.String()
}

func (t UnionType) ElementType() Type {
	var elType Type = VoidType{}
	for _, disjunct := range t {
		elType = Union(elType, disjunct.ElementType())
	}
	return elType
}
