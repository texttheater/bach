package types

type Void struct {
}

func (t Void) Subsumes(u Type) bool {
	switch u.(type) {
	case Void:
		return true
	default:
		return false
	}
}

func (t Void) Bind(u Type, bindings map[string]Type) bool {
	switch u.(type) {
	case Void:
		return true
	default:
		return false
	}
}

func (t Void) Instantiate(bindings map[string]Type) Type {
	return t
}

func (t Void) Partition(u Type) (Type, Type) {
	return t, t
}

func (t Void) String() string {
	return "Void"
}

func (t Void) ElementType() Type {
	panic("Void is not a sequence type")
}
