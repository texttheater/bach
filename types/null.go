package types

var NullType = nullType{}

type nullType struct {
}

func (t nullType) Subsumes(other Type) bool {
	if VoidType.Subsumes(other) {
		return true
	}
	_, ok := other.(nullType)
	return ok
}

func (t nullType) String() string {
	return "Null"
}

func (t nullType) ElementType() Type {
	panic("Null is not a sequence type")
}
