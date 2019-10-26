package types

type NumType struct {
}

func (t NumType) Subsumes(u Type) bool {
	switch u.(type) {
	case VoidType:
		return true
	case NumType:
		return true
	default:
		return false
	}
}

func (t NumType) Bind(u Type, bindings map[string]Type) bool {
	switch u.(type) {
	case VoidType:
		return true
	case NumType:
		return true
	default:
		return false
	}
}

func (t NumType) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case VoidType:
		return u, t
	case NumType:
		return u, VoidType{}
	case UnionType:
		return u.inversePartition(t)
	case AnyType:
		return t, VoidType{}
	default:
		return VoidType{}, t
	}
}

func (t NumType) String() string {
	return "Num"
}

func (t NumType) ElementType() Type {
	panic("Num is not a sequence type")
}
