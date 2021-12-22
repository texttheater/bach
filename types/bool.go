package types

type Bool struct {
}

func (t Bool) Subsumes(u Type) bool {
	switch u.(type) {
	case Void:
		return true
	case Bool:
		return true
	default:
		return false
	}
}

func (t Bool) Bind(u Type, bindings map[string]Type) bool {
	switch u.(type) {
	case Void:
		return true
	case Bool:
		return true
	default:
		return false
	}
}

func (t Bool) Instantiate(bindings map[string]Type) Type {
	return t
}

func (t Bool) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case Void:
		return u, t
	case Bool:
		return u, Void{}
	case Union:
		return u.inversePartition(t)
	case Any:
		return t, Void{}
	default:
		return Void{}, t
	}
}

func (t Bool) String() string {
	return "Bool"
}

func (t Bool) ElementType() Type {
	panic("Bool is not a sequence type")
}
