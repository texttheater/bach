package types

import (
	"fmt"
)

func SeqType(elementType Type) Type {
	return seqType{elementType}
}

type seqType struct {
	elementType Type
}

// AnySeqType represents Seq<Any>, the type of sequences with any type of
// element. It is provided as a package variable for convenience.
var AnySeqType = SeqType(AnyType())

func (t seqType) Subsumes(u Type) bool {
	if VoidType.Subsumes(u) {
		return true
	}
	if ArrType(t.elementType).Subsumes(u) {
		return true
	}
	switch u := u.(type) {
	case seqType:
		return t.elementType.Subsumes(u.elementType)
	case unionType:
		for _, disjunct := range u {
			if !t.Subsumes(disjunct) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func (t seqType) String() string {
	return fmt.Sprintf("Seq<%v>", t.elementType)
}

func (t seqType) ElementType() Type {
	return t.elementType
}
