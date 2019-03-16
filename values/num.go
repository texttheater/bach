package values

import (
	"fmt"
	"strconv"
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
