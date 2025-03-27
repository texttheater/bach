package errors

import (
	"fmt"
)

type ErrorKind int

const (
	SyntaxKind ErrorKind = iota
	TypeKind
	ValueKind
	UnknownKind
)

func (kind ErrorKind) String() string {
	switch kind {
	case SyntaxKind:
		return "Syntax error"
	case TypeKind:
		return "Type error"
	case ValueKind:
		return "Value error"
	default:
		return "Unknown error"
	}
}

func ParseKind(s string) (ErrorKind, error) {
	switch s {
	case "Syntax":
		return SyntaxKind, nil
	case "Type":
		return TypeKind, nil
	case "Value":
		return ValueKind, nil
	case "UnknownKind":
		return UnknownKind, nil
	default:
		return 0, fmt.Errorf("invalid error kind")
	}
}
