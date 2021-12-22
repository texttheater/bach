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
	case types.Reader:
		return true, nil
	case types.Union:
		return inhabits(v, t, stack)
	case types.Any:
		return true, nil
	case types.Var:
		return stack.Inhabits(v, t)
	default:
		return false, nil
	}
}

func (v ReaderValue) Equal(w Value) (bool, error) {
	return v == w, nil
}
