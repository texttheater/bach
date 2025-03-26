package types

type VoidType struct {
}

func (t VoidType) Subsumes(u Type) bool {
	switch u.(type) {
	case VoidType:
		return true
	default:
		return false
	}
}

func (t VoidType) Bind(u Type, bindings map[string]Type) bool {
	switch u.(type) {
	case VoidType:
		return true
	default:
		return false
	}
}

func (t VoidType) Instantiate(bindings map[string]Type) Type {
	return t
}

func (t VoidType) Partition(u Type) (Type, Type) {
	return t, t
}

func (t VoidType) String() string {
	return "Void"
}

func (t VoidType) ElementType() Type {
	panic("Void is not a sequence type")
}
