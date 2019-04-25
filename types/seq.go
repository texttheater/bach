package types

import (
	"fmt"
)

type SeqType struct {
	ElType Type
}

// AnySeqType represents Seq<Any>, the type of sequences with any type of
// element. It is provided as a package variable for convenience.
var AnySeqType Type = &SeqType{AnyType{}}

func (t *SeqType) Subsumes(u Type) bool {
	if (VoidType{}).Subsumes(u) {
		return true
	}
	if (&ArrType{t.ElType}).Subsumes(u) {
		return true
	}
	switch u := u.(type) {
	case *SeqType:
		return t.ElType.Subsumes(u.ElType)
	case UnionType:
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

func (t *SeqType) String() string {
	return fmt.Sprintf("Seq<%v>", t.ElType)
}

func (t *SeqType) ElementType() Type {
	return t.ElType
}
