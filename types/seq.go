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
	switch u := u.(type) {
	case VoidType:
		return true
	case *NearrType:
		if !t.ElType.Subsumes(u.HeadType) {
			return false
		}
		return t.Subsumes(u.TailType)
	case *ArrType:
		return t.ElType.Subsumes(u.ElType)
	case *SeqType:
		return t.ElType.Subsumes(u.ElType)
	case UnionType:
		return u.inverseSubsumes(t)
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
