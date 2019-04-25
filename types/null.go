package types

type NullType struct {
}

func (t NullType) Subsumes(other Type) bool {
	if (VoidType{}).Subsumes(other) {
		return true
	}
	_, ok := other.(NullType)
	return ok
}

func (t NullType) String() string {
	return "Null"
}

func (t NullType) ElementType() Type {
	panic("Null is not a sequence type")
}
