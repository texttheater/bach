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
var AnySeqType = SeqType(AnyType)

func (t seqType) Subsumes(other Type) bool {
	if VoidType.Subsumes(other) {
		return true
	}
	otherSeqType, ok := other.(seqType)
	if ok {
		return t.elementType.Subsumes(otherSeqType.elementType)
	}
	if ArrType(t.elementType).Subsumes(other) {
		return true
	}
	return false
}

func (t seqType) String() string {
	return fmt.Sprintf("Seq<%v>", t.elementType)
}

func (t seqType) ElementType() Type {
	return t.elementType
}
