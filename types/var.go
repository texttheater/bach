package types

type TypeVariable struct {
	Name string
}

func (t TypeVariable) Subsumes(u Type) bool {
	switch u := u.(type) {
	case VoidType:
		return true
	case TypeVariable:
		return t.Name == u.Name
	default:
		return false
	}
}

func (t TypeVariable) Partition(u Type) (Type, Type) {
	panic("cannot partition a type variable")
}

func (t TypeVariable) String() string {
	return t.Name
}

func (t TypeVariable) ElementType() Type {
	panic("type variable " + t.Name + " is not a sequence type")
}
