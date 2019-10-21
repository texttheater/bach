package values

import (
	"bytes"

	"github.com/texttheater/bach/types"
)

func NewArrValue(elements []Value) *ArrValue {
	var v *ArrValue = nil
	for i := len(elements) - 1; i >= 0; i-- {
		v = &ArrValue{
			Head: elements[i],
			Tail: v,
		}
	}
	return v
}

type ArrValue struct {
	Head Value
	Tail *ArrValue
}

func (v *ArrValue) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("[")
	if v != nil {
		buffer.WriteString(v.Head.String())
		next := v.Tail
		for next != nil {
			buffer.WriteString(", ")
			buffer.WriteString(next.Head.String())
			next = next.Tail
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
		for v != nil {
			channel <- v.Head
			v = v.Tail
		}
		close(channel)
	}()
	return channel
}

func (v *ArrValue) Inhabits(t types.Type) bool {
	switch t := t.(type) {
	case types.TupType:
		for _, elType := range t {
			if v == nil {
				return false
			}
			if !v.Head.Inhabits(elType) {
				return false
			}
			v = v.Tail
		}
		if v != nil {
			return false
		}
		return true
	case *types.ArrType:
		if (types.AnyType{}).Subsumes(t.ElType) {
			return true
		}
		for v != nil {
			if !v.Head.Inhabits(t.ElType) {
				return false
			}
			v = v.Tail
		}
		return true
	case *types.SeqType:
		if (types.AnyType{}).Subsumes(t.ElType) {
			return true
		}
		for v != nil {
			if !v.Head.Inhabits(t.ElType) {
				return false
			}
			v = v.Tail
		}
		return true
	case types.UnionType:
		return inhabits(v, t)
	case types.AnyType:
		return true
	default:
		return false

	}
}
