package errors

import (
	"encoding/json"
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

func (kind *ErrorKind) UnmarshalJSON(data []byte) error {
	var v string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	switch v {
	case "Syntax":
		*kind = SyntaxKind
	case "Type":
		*kind = TypeKind
	case "Value":
		*kind = ValueKind
	case "UnknownKind":
		*kind = UnknownKind
	default:
		return fmt.Errorf("invalid error kind")
	}
	return nil
}
