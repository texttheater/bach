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
		Ch: channel,
	}
}

type ArrValue struct {
	Ch <-chan Value
	Hd Value
	Tl *ArrValue
}

func (v *ArrValue) Eval() {
	if v.Ch == nil {
		return
	}
	v.Hd = <-v.Ch
	v.Tl = &ArrValue{
		Ch: v.Ch,
	}
	v.Ch = nil
}

func (v *ArrValue) IsEmpty() bool {
	v.Eval()
	return v.Hd == nil
}

func (v *ArrValue) Length() int {
	length := 0
	for !v.IsEmpty() {
		length += 1
		v = v.Tl
	}
	return length
}

func (v *ArrValue) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("[")
	if !v.IsEmpty() {
		buffer.WriteString(v.Hd.String())
		v = v.Tl
		for !v.IsEmpty() {
			buffer.WriteString(", ")
			buffer.WriteString(v.Hd.String())
			v = v.Tl
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
			channel <- v.Hd
			v = v.Tl
		}
		close(channel)
	}()
	return channel
}

func (v *ArrValue) Inhabits(t types.Type, stack *BindingStack) bool {
	switch t := t.(type) {
	case types.TupType:
		for i := range t {
			if v.IsEmpty() {
				return false
			}
			if !v.Hd.Inhabits(t[i], stack) {
				return false
			}
			v = v.Tl
		}
		if !v.IsEmpty() {
			return false
		}
		return true
	case *types.ArrType:
		if (types.AnyType{}).Subsumes(t.ElType) {
			return true
		}
		for !v.IsEmpty() {
			if !v.Hd.Inhabits(t.ElType, stack) {
				return false
			}
			v = v.Tl
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
		if !v.Hd.Equal(w.Hd) {
			return false
		}
		return v.Tl.Equal(w.Tl)
	default:
		return false
	}
}
