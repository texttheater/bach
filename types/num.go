package types

type NumType struct {
}

func (t NumType) Subsumes(other Type) bool {
	if (VoidType{}).Subsumes(other) {
		return true
	}
	_, ok := other.(NumType)
	return ok
}

func (t NumType) String() string {
	return "Num"
}

func (t NumType) ElementType() Type {
	panic("Num is not a sequence type")
}
