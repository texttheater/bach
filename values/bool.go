package values

import (
	"fmt"
	"strconv"
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
