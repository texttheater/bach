package types

type StrType struct {
}

func (t StrType) Subsumes(u Type) bool {
	switch u.(type) {
	case VoidType:
		return true
	case StrType:
		return true
	default:
		return false
	}
}

func (t StrType) String() string {
	return "Str"
}

func (t StrType) ElementType() Type {
	panic("Str is not a sequence type")
}
