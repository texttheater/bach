package values

import (
	"bytes"

	"github.com/texttheater/bach/types"
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

func (v ArrValue) Inhabits(t types.Type, stack *BindingStack) bool {
	switch t := t.(type) {
	case types.TupType:
		if len(v) != len(t) {
			return false
		}
		for i := range v {
			if !v[i].Inhabits(t[i], stack) {
				return false
			}
		}
		return true
	case *types.ArrType:
		if (types.AnyType{}).Subsumes(t.ElType) {
			return true
		}
		for _, e := range v {
			if !e.Inhabits(t.ElType, stack) {
				return false
			}
		}
		return true
	case *types.SeqType:
		if (types.AnyType{}).Subsumes(t.ElType) {
			return true
		}
		for _, e := range v {
			if !e.Inhabits(t.ElType, stack) {
				return false
			}
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
