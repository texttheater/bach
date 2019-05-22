package values

import (
	"fmt"
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

func (v ReaderValue) Iter() <-chan Value {
	panic(fmt.Sprintf("%s is not a sequence", v))
}

func (v ReaderValue) Inhabits(t types.Type) bool {
	switch t := t.(type) {
	case types.ReaderType:
		return true
	case types.UnionType:
		return inhabits(v, t)
	case types.AnyType:
		return true
	default:
		return false
	}
}
