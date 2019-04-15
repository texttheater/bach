package values

import (
	"fmt"
)

type NullValue struct {
}

func (v NullValue) String() string {
	return "null"
}

func (v NullValue) Out() string {
	return v.String()
}

func (v NullValue) Iter() <-chan Value {
	panic(fmt.Sprintf("%s is not a sequence", v))
}
