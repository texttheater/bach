package types

type BoolType struct {
}

func (t BoolType) Subsumes(u Type) bool {
	switch u.(type) {
	case VoidType:
		return true
	case BoolType:
		return true
	default:
		return false
	}
}

func (t BoolType) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case VoidType:
		return u, t
	case BoolType:
		return u, VoidType{}
	case UnionType:
		return u.inversePartition(t)
	case AnyType:
		return t, VoidType{}
	default:
		return VoidType{}, t
	}
}

func (t BoolType) String() string {
	return "Bool"
}

func (t BoolType) ElementType() Type {
	panic("Bool is not a sequence type")
}
