package types

type StrType struct {
}

func (t StrType) Subsumes(u Type) bool {
	switch u.(type) {
	case VoidType:
		return true
	case StrType:
		return true
	default:
		return false
	}
}

func (t StrType) Bind(u Type, bindings map[string]Type) bool {
	switch u.(type) {
	case VoidType:
		return true
	case StrType:
		return true
	default:
		return false
	}
}

func (t StrType) Instantiate(bindings map[string]Type) Type {
	return t
}

func (t StrType) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case VoidType:
		return u, t
	case StrType:
		return u, VoidType{}
	case UnionType:
		return u.inversePartition(t)
	case AnyType:
		return t, VoidType{}
	default:
		return VoidType{}, t
	}
}

func (t StrType) String() string {
	return "Str"
}

func (t StrType) ElementType() Type {
	panic("Str is not a sequence type")
}
