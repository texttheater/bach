package types

var StrType = strType{}

type strType struct {
}

func (t strType) Subsumes(other Type) bool {
	if VoidType.Subsumes(other) {
		return true
	}
	_, ok := other.(strType)
	return ok
}

func (t strType) String() string {
	return "Str"
}

func (t strType) ElementType() Type {
	panic("Str is not a sequence type")
}
