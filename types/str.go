package types

type Str struct {
}

func (t Str) Subsumes(u Type) bool {
	switch u.(type) {
	case Void:
		return true
	case Str:
		return true
	default:
		return false
	}
}

func (t Str) Bind(u Type, bindings map[string]Type) bool {
	switch u.(type) {
	case Void:
		return true
	case Str:
		return true
	default:
		return false
	}
}

func (t Str) Instantiate(bindings map[string]Type) Type {
	return t
}

func (t Str) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case Void:
		return u, t
	case Str:
		return u, Void{}
	case Union:
		return u.inversePartition(t)
	case Any:
		return t, Void{}
	default:
		return Void{}, t
	}
}

func (t Str) String() string {
	return "Str"
}

func (t Str) ElementType() Type {
	panic("Str is not a sequence type")
}
