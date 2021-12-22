package types

type Reader struct {
}

func (t Reader) Subsumes(u Type) bool {
	switch u.(type) {
	case Void:
		return true
	case Reader:
		return true
	default:
		return false
	}
}

func (t Reader) Bind(u Type, bindings map[string]Type) bool {
	switch u.(type) {
	case Void:
		return true
	case Reader:
		return true
	default:
		return false
	}
}

func (t Reader) Instantiate(bindings map[string]Type) Type {
	return t
}

func (t Reader) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case Void:
		return u, t
	case Reader:
		return u, Void{}
	case Union:
		return u.inversePartition(t)
	case Any:
		return t, Void{}
	default:
		return Void{}, t
	}
}

func (t Reader) String() string {
	return "Reader"
}

func (t Reader) ElementType() Type {
	panic("Reader is not a sequence type")
}
