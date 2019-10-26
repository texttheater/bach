package values

import (
	"fmt"
	"strconv"

	"github.com/texttheater/bach/types"
)

type NumValue float64

func (v NumValue) String() string {
	return strconv.FormatFloat(float64(v), 'f', -1, 64)
}

func (v NumValue) Out() string {
	return v.String()
}

func (v NumValue) Iter() <-chan Value {
	panic(fmt.Sprintf("%s is not a sequence", v))
}

func (v NumValue) Inhabits(t types.Type, stack *BindingStack) bool {
	switch t := t.(type) {
	case types.NumType:
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
