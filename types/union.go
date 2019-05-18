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

// TODO make this private?

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
