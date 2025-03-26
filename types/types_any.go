package types

type Any struct {
}

func (t Any) Subsumes(u Type) bool {
	return true
}

func (t Any) Bind(u Type, bindings map[string]Type) bool {
	return true
}

func (t Any) Instantiate(bindings map[string]Type) Type {
	return t
}

func (t Any) Partition(u Type) (Type, Type) {
	switch u.(type) {
	case Any:
		return u, Void{}
	default:
		return u, t
	}
}

func (t Any) String() string {
	return "Any"
}

func (t Any) ElementType() Type {
	panic("Any is not a sequence type")
}
