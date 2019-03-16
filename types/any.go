package types

var AnyType = anyType{}

type anyType struct {
}

func (t anyType) Subsumes(other Type) bool {
	return true
}

func (t anyType) String() string {
	return "Any"
}

func (t anyType) ElementType() Type {
	panic("Any is not a sequence type")
}
