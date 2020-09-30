package states

import (
	"io"

	"github.com/texttheater/bach/types"
)

type ReaderValue struct {
	Reader io.Reader
}

func (v ReaderValue) Repr() (string, error) {
	return "<reader>", nil
}

func (v ReaderValue) Str() (string, error) {
	return v.Repr()
}

func (v ReaderValue) Inhabits(t types.Type, stack *BindingStack) (bool, error) {
	switch t := t.(type) {
	case types.ReaderType:
		return true, nil
	case types.UnionType:
		return inhabits(v, t, stack)
	case types.AnyType:
		return true, nil
	case types.TypeVariable:
		return stack.Inhabits(v, t)
	default:
		return false, nil
	}
}

func (v ReaderValue) Equal(w Value) (bool, error) {
	return v == w, nil
}
