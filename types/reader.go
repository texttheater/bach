package types

type ReaderType struct {
}

func (t ReaderType) Subsumes(u Type) bool {
	switch u.(type) {
	case VoidType:
		return true
	case ReaderType:
		return true
	default:
		return false
	}
}

func (t ReaderType) Bind(u Type, bindings map[string]Type) bool {
	switch u.(type) {
	case VoidType:
		return true
	case ReaderType:
		return true
	default:
		return false
	}
}

func (t ReaderType) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case VoidType:
		return u, t
	case ReaderType:
		return u, VoidType{}
	case UnionType:
		return u.inversePartition(t)
	case AnyType:
		return t, VoidType{}
	default:
		return VoidType{}, t
	}
}

func (t ReaderType) String() string {
	return "Reader"
}

func (t ReaderType) ElementType() Type {
	panic("Reader is not a sequence type")
}
