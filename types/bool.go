package types

type BoolType struct {
}

func (t BoolType) Subsumes(other Type) bool {
	if (VoidType{}).Subsumes(other) {
		return true
	}
	_, ok := other.(BoolType)
	return ok
}

func (t BoolType) String() string {
	return "Bool"
}

func (t BoolType) ElementType() Type {
	panic("Bool is not a sequence type")
}
