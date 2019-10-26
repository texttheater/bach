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

func (t TypeVariable) Bind(u Type, bindings map[string]Type) bool {
	instType, ok := bindings[t.Name]
	if !ok {
		instType = AnyType{}
	}
	var newInstType Type
	// pick the more specific type
	if instType.Subsumes(u) {
		newInstType = u
	} else if u.Subsumes(instType) {
		newInstType = instType
	} else {
		return false
	}
	bindings[t.Name] = newInstType
	return true
}

func (t TypeVariable) Instantiate(bindings map[string]Type) Type {
	instType, ok := bindings[t.Name]
	if !ok {
		return t
	}
	return instType
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
