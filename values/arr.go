package values

import (
	"bytes"
)

type ArrValue []Value

func (v ArrValue) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("[")
	if len(v) > 0 {
		buffer.WriteString(v[0].String())
		for _, elValue := range v[1:] {
			buffer.WriteString(", ")
			buffer.WriteString(elValue.String())
		}
	}
	buffer.WriteString("]")
	return buffer.String()
}

func (v ArrValue) Out() string {
	return v.String()
}

func (v ArrValue) Iter() <-chan Value {
	channel := make(chan Value)
	go func() {
		for _, el := range v {
			channel <- el
		}
		close(channel)
	}()
	return channel
}
