package types

type VoidType struct {
}

func (t VoidType) Subsumes(other Type) bool {
	_, ok := other.(VoidType)
	return ok
}

func (t VoidType) String() string {
	return "Void"
}

func (t VoidType) ElementType() Type {
	panic("Void is not a sequence type")
}
