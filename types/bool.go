package types

var BoolType = boolType{}

type boolType struct {
}

func (t boolType) Subsumes(other Type) bool {
	if VoidType.Subsumes(other) {
		return true
	}
	_, ok := other.(boolType)
	return ok
}

func (t boolType) String() string {
	return "Bool"
}

func (t boolType) ElementType() Type {
	panic("Bool is not a sequence type")
}
