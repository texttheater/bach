package values

import (
	"bytes"
	"fmt"
)

type ObjValue map[string]Value

func (v ObjValue) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("{")
	firstWritten := false
	for k, w := range v {
		if firstWritten {
			buffer.WriteString(", ")
		}
		buffer.WriteString(k)
		buffer.WriteString(": ")
		buffer.WriteString(w.String())
		firstWritten = true
	}
	buffer.WriteString("}")
	return buffer.String()
}

func (v ObjValue) Out() string {
	return v.String()
}

func (v ObjValue) Iter() <-chan Value {
	panic(fmt.Sprintf("%s is not a sequence", v))
}
