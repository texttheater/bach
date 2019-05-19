package types

type AnyType struct {
}

func (t AnyType) Subsumes(u Type) bool {
	return true
}

func (t AnyType) String() string {
	return "Any"
}

func (t AnyType) ElementType() Type {
	panic("Any is not a sequence type")
}
