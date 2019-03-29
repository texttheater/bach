package values

import (
	"fmt"
	"io"
)

type ReaderValue struct {
	Reader io.Reader
}

func (v ReaderValue) String() string {
	return fmt.Sprintf("%s", v.Reader)
}

func (v ReaderValue) Out() string {
	return v.String()
}

func (v ReaderValue) Iter() <-chan Value {
	panic(fmt.Sprintf("%s is not a sequence", v))
}
