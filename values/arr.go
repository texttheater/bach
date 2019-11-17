package values

import (
	"bytes"

	"github.com/texttheater/bach/types"
)

func NewArrValue(elements []Value) *ArrValue {
	i := 0
	var next func() (Value, *ArrValue, error)
	next = func() (Value, *ArrValue, error) {
		if i >= len(elements) {
			return nil, nil, nil
		}
		head := elements[i]
		i++
		return head, &ArrValue{
			Func: next,
		}, nil
	}
	return &ArrValue{
		Func: next,
	}
}

type ArrValue struct {
	Func func() (Value, *ArrValue, error)
	Head Value
	Tail *ArrValue
}

func (v *ArrValue) Eval() error {
	if v.Func == nil {
		return nil
	}
	var err error
	v.Head, v.Tail, err = v.Func()
	v.Func = nil
	return err
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
