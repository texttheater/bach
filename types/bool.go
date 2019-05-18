package types

type BoolType struct {
}

func (t BoolType) Subsumes(u Type) bool {
	switch u.(type) {
	case VoidType:
		return true
	case BoolType:
		return true
	default:
		return false
	}
}

func (t BoolType) String() string {
	return "Bool"
}

func (t BoolType) ElementType() Type {
	panic("Bool is not a sequence type")
}
