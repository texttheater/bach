package types

type Type interface {
	Subsumes(Type) bool
	Bind(Type, map[string]Type) bool
	Partition(Type) (Type, Type)
	String() string
	ElementType() Type
}
