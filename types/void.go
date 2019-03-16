package types

var VoidType = voidType{}

type voidType struct {
}

func (t voidType) Subsumes(other Type) bool {
	_, ok := other.(voidType)
	return ok
}

func (t voidType) String() string {
	return "Void"
}

func (t voidType) ElementType() Type {
	panic("Void is not a sequence type")
}
