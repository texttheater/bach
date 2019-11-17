package values

import (
	"io"

	"github.com/texttheater/bach/types"
)

type ReaderValue struct {
	Reader io.Reader
}

func (v ReaderValue) String() string {
	return "<reader>"
}

func (v ReaderValue) Out() string {
	return v.String()
}

func (v ReaderValue) Inhabits(t types.Type, stack *BindingStack) bool {
	switch t := t.(type) {
	case types.ReaderType:
		return true
	case types.UnionType:
		return inhabits(v, t, stack)
	case types.AnyType:
		return true
	case types.TypeVariable:
		return stack.Inhabits(v, t)
	default:
		return false
	}
}

func (v ReaderValue) Equal(w Value) bool {
	return v == w
}
