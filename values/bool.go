package values

import (
	"fmt"
	"strconv"

	"github.com/texttheater/bach/types"
)

type BoolValue bool

func (v BoolValue) String() string {
	return strconv.FormatBool(bool(v))
}

func (v BoolValue) Out() string {
	return v.String()
}

func (v BoolValue) Iter() <-chan Value {
	panic(fmt.Sprintf("%s is not a sequence", v))
}

func (v BoolValue) Inhabits(t types.Type) bool {
	switch t := t.(type) {
	case types.BoolType:
		return true
	case types.UnionType:
		return inhabits(v, t)
	case types.AnyType:
		return true
	default:
		return false
	}
}
