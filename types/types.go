package types

type Type interface {
	Subsumes(Type) bool
	String() string
	ElementType() Type
}
