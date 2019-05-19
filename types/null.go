package types

type NullType struct {
}

func (t NullType) Subsumes(u Type) bool {
	switch u.(type) {
	case VoidType:
		return true
	case NullType:
		return true
	default:
		return false
	}
}

func (t NullType) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case VoidType:
		return u, t
	case NullType:
		return u, VoidType{}
	case UnionType:
		return u.inversePartition(t)
	case AnyType:
		return t, VoidType{}
	default:
		return VoidType{}, t
	}
}

func (t NullType) String() string {
	return "Null"
}

func (t NullType) ElementType() Type {
	panic("Null is not a sequence type")
}
