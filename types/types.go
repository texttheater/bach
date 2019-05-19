package types

type Type interface {
	Subsumes(Type) bool
	Partition(Type) (Type, Type)
	String() string
	ElementType() Type
}
