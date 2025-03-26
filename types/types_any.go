package types

type AnyType struct {
}

func (t AnyType) Subsumes(u Type) bool {
	return true
}

func (t AnyType) Bind(u Type, bindings map[string]Type) bool {
	return true
}

func (t AnyType) Instantiate(bindings map[string]Type) Type {
	return t
}

func (t AnyType) Partition(u Type) (Type, Type) {
	switch u.(type) {
	case AnyType:
		return u, VoidType{}
	default:
		return u, t
	}
}

func (t AnyType) String() string {
	return "Any"
}

func (t AnyType) ElementType() Type {
	panic("Any is not a sequence type")
}
