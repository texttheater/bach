package types

type StrType struct {
}

func (t StrType) Subsumes(other Type) bool {
	if (VoidType{}).Subsumes(other) {
		return true
	}
	_, ok := other.(StrType)
	return ok
}

func (t StrType) String() string {
	return "Str"
}

func (t StrType) ElementType() Type {
	panic("Str is not a sequence type")
}
