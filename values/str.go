package values

import (
	"fmt"
)

type StrValue string

func (v StrValue) String() string {
	return fmt.Sprintf("%q", string(v))
}

func (v StrValue) Out() string {
	return string(v)
}

func (v StrValue) Iter() <-chan Value {
	panic(fmt.Sprintf("%s is not a sequence", v))
}
