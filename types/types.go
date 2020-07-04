package types

type Type interface {
	Subsumes(Type) bool
	Bind(Type, map[string]Type) bool
	Instantiate(map[string]Type) Type
	Partition(Type) (Type, Type)
	String() string
	ElementType() Type
}

func Equivalent(a Type, b Type) bool {
	return a.Subsumes(b) && b.Subsumes(a)
}
