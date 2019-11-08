package values

import (
	"bytes"

	"github.com/texttheater/bach/types"
)

func NewArrValue(elements []Value) *ArrValue {
	channel := make(chan Value)
	go func() {
		for _, el := range elements {
			channel <- el
		}
		close(channel)
	}()
	return &ArrValue{
		Channel: channel,
	}
}

type ArrValue struct {
	Channel <-chan Value
	Head    Value
	Tail    *ArrValue
}

func (v *ArrValue) Eval() {
	if v.Channel == nil {
		return
	}
	v.Head = <-v.Channel
	v.Tail = &ArrValue{
		Channel: v.Channel,
	}
	v.Channel = nil
}

func (v *ArrValue) IsEmpty() bool {
	v.Eval()
	return v.Head == nil
}

func (v *ArrValue) Length() int {
	length := 0
	for !v.IsEmpty() {
		length += 1
		v = v.Tail
	}
	return length
}

func (v *ArrValue) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("[")
	if !v.IsEmpty() {
		buffer.WriteString(v.Head.String())
		v = v.Tail
		for !v.IsEmpty() {
			buffer.WriteString(", ")
			buffer.WriteString(v.Head.String())
			v = v.Tail
		}
	}
	buffer.WriteString("]")
	return buffer.String()
}

func (v *ArrValue) Out() string {
	return v.String()
}

func (v *ArrValue) Iter() <-chan Value {
	channel := make(chan Value)
	go func() {
		for !v.IsEmpty() {
			channel <- v.Head
			v = v.Tail
		}
		close(channel)
	}()
	return channel
}

func (v *ArrValue) Inhabits(t types.Type, stack *BindingStack) bool {
	switch t := t.(type) {
	case *types.NearrType:
		if v.IsEmpty() {
			return false
		}
		if !v.Head.Inhabits(t.HeadType, stack) {
			return false
		}
		return v.Tail.Inhabits(t.TailType, stack)
	case *types.ArrType:
		if (types.AnyType{}).Subsumes(t.ElType) {
			return true
		}
		for !v.IsEmpty() {
			if !v.Head.Inhabits(t.ElType, stack) {
				return false
			}
			v = v.Tail
		}
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

func (v *ArrValue) Equal(w Value) bool {
	switch w := w.(type) {
	case *ArrValue:
		if v.IsEmpty() {
			return w.IsEmpty()
		}
		if w.IsEmpty() {
			return false
		}
		if !v.Head.Equal(w.Head) {
			return false
		}
		return v.Tail.Equal(w.Tail)
	default:
		return false
	}
}
