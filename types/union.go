package types

import (
	"bytes"
	"fmt"
)

func Union(a Type, b Type) Type {
	if a.Subsumes(b) {
		return a
	}
	if b.Subsumes(a) {
		return b
	}
	aDisjuncts := disjuncts(a)
	bDisjuncts := disjuncts(b)
	var disjuncts []Type
	added := make([]bool, len(bDisjuncts))
aLoop:
	for _, aDisjunct := range aDisjuncts {
		for i, bDisjunct := range bDisjuncts {
			if added[i] {
				continue
			}
			if bDisjunct.Subsumes(aDisjunct) {
				disjuncts = append(disjuncts, bDisjunct)
				added[i] = true
				continue aLoop
			}
		}
		disjuncts = append(disjuncts, aDisjunct)
	}
	for i, a := range added {
		if !a {
			disjuncts = append(disjuncts, bDisjuncts[i])
		}
	}
	return unionType(disjuncts)
}

func disjuncts(t Type) []Type {
	tUnion, ok := t.(unionType)
	if ok {
		return tUnion
	}
	return []Type{t}
}

type unionType []Type

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
