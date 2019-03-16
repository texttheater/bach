package types

import (
	"bytes"
	"fmt"
)

func Union(a Type, b Type) Type {
	aUnion, ok := a.(unionType)
	if ok {
		bUnion, ok := b.(unionType)
		if ok {
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
	return unionType([]Type{a, b})
}

type unionType []Type

func typeAppend(t unionType, u Type) unionType {
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

func (t unionType) Subsumes(u Type) bool {
	uUnion, ok := u.(unionType)
	if ok {
	uLoop:
		for _, uDisjunct := range uUnion {
			// find a subsumer for uDisjunct among t
			for _, tDisjunct := range t {
				if tDisjunct.Subsumes(uDisjunct) {
					// subsumer found, check next uDisjunct
					continue uLoop
				}
			}
			// no subsumer found
			return false
		}
		// all uDisjuncts checked
		return true
	}
	// find a subsumer for u
	for _, tDisjunct := range t {
		if tDisjunct.Subsumes(u) {
			return true
		}
	}
	// no subsumer found
	return false
}

func (t unionType) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString(fmt.Sprintf("%s", t[0]))
	for _, disjunct := range t[1:] {
		buffer.WriteString("|")
		buffer.WriteString(fmt.Sprintf("%s", disjunct))
	}
	return buffer.String()
}

func (t unionType) ElementType() Type {
	panic(t.String() + " is not a sequence type")
}
