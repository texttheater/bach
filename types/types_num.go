package types

type Num struct {
}

func (t Num) Subsumes(u Type) bool {
	switch u.(type) {
	case Void:
		return true
	case Num:
		return true
	default:
		return false
	}
}

func (t Num) Bind(u Type, bindings map[string]Type) bool {
	switch u.(type) {
	case Void:
		return true
	case Num:
		return true
	default:
		return false
	}
}

func (t Num) Instantiate(bindings map[string]Type) Type {
	return t
}

func (t Num) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case Void:
		return u, t
	case Num:
		return u, Void{}
	case Union:
		return u.inversePartition(t)
	case Any:
		return t, Void{}
	default:
		return Void{}, t
	}
}

func (t Num) String() string {
	return "Num"
}

func (t Num) ElementType() Type {
	panic("Num is not a sequence type")
}
