package types

import (
	"bytes"
	"fmt"
)

func Disjoin(a Type, b Type) Type {
	aDisj, ok := a.(disjunctiveType)
	if ok {
		return aDisj.disjoin(b)
	}
	bDisj, ok := b.(disjunctiveType)
	if ok {
		return bDisj.disjoin(a)
	}
	if a.Subsumes(b) {
		return a
	}
	if b.Subsumes(a) {
		return b
	}
	return disjunctiveType{[]Type{a, b}}
}

type disjunctiveType struct {
	disjuncts []Type
}

func (t disjunctiveType) Subsumes(other Type) bool {
	otherDisj, ok := other.(disjunctiveType)
	if ok {
		return t.subsumesDisj(otherDisj)
	}
	return t.subsumesNonDisj(other)
}

func (t disjunctiveType) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString(fmt.Sprintf("%s", t.disjuncts[0]))
	for _, disjunct := range t.disjuncts[1:] {
		buffer.WriteString("|")
		buffer.WriteString(fmt.Sprintf("%s", disjunct))
	}
	return buffer.String()
}

func (t disjunctiveType) ElementType() Type {
	panic(t.String() + " is not a sequence type")
}

func (t disjunctiveType) subsumesDisj(other disjunctiveType) bool {
	for _, disjunct := range other.disjuncts {
		if !t.subsumesNonDisj(disjunct) {
			return false
		}
	}
	return true
}

func (t disjunctiveType) subsumesNonDisj(other Type) bool {
	for _, disjunct := range t.disjuncts {
		if disjunct.Subsumes(other) {
			return true
		}
	}
	return false
}

func (t disjunctiveType) disjoin(other Type) Type {
	otherDisj, ok := other.(disjunctiveType)
	if ok {
		return t.disjoinDisj(otherDisj)
	}
	return t.disjoinNonDisj(other)
}

func (t disjunctiveType) disjoinDisj(other disjunctiveType) Type {
	result := t
	for _, disjunct := range other.disjuncts {
		result = result.disjoinNonDisj(disjunct)
	}
	return result
}

func (t disjunctiveType) disjoinNonDisj(other Type) disjunctiveType {
	for _, disjunct := range t.disjuncts {
		if disjunct.Subsumes(other) {
			return t
		}
	}
	newDisjuncts := make([]Type, len(t.disjuncts)+1)
	for i, disjunct := range t.disjuncts {
		if !other.Subsumes(disjunct) {
			newDisjuncts[i] = disjunct
		}
	}
	newDisjuncts[len(t.disjuncts)] = other
	return disjunctiveType{newDisjuncts}
}
