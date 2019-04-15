package types

var NumType = numType{}

type numType struct {
}

func (t numType) Subsumes(other Type) bool {
	if VoidType.Subsumes(other) {
		return true
	}
	_, ok := other.(numType)
	return ok
}

func (t numType) String() string {
	return "Num"
}

func (t numType) ElementType() Type {
	panic("Num is not a sequence type")
}
