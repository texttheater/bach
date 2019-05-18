package types

type NumType struct {
}

func (t NumType) Subsumes(u Type) bool {
	switch u.(type) {
	case VoidType:
		return true
	case NumType:
		return true
	default:
		return false
	}
}

func (t NumType) String() string {
	return "Num"
}

func (t NumType) ElementType() Type {
	panic("Num is not a sequence type")
}
