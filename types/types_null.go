package types

type Null struct {
}

func (t Null) Subsumes(u Type) bool {
	switch u.(type) {
	case Void:
		return true
	case Null:
		return true
	default:
		return false
	}
}

func (t Null) Bind(u Type, bindings map[string]Type) bool {
	switch u.(type) {
	case Void:
		return true
	case Null:
		return true
	default:
		return false
	}
}

func (t Null) Instantiate(bindings map[string]Type) Type {
	return t
}

func (t Null) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case Void:
		return u, t
	case Null:
		return u, Void{}
	case Union:
		return u.inversePartition(t)
	case Any:
		return t, Void{}
	default:
		return Void{}, t
	}
}

func (t Null) String() string {
	return "Null"
}

func (t Null) ElementType() Type {
	panic("Null is not a sequence type")
}
