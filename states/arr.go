package states

import (
	"bytes"

	"github.com/texttheater/bach/types"
)

func NewArrValue(elements []Value) *ArrValue {
	var arrFrom func(i int) *ArrValue
	arrFrom = func(i int) *ArrValue {
		if i == len(elements) {
			return nil
		}
		return &ArrValue{
			Head: elements[i],
			Tail: &Thunk{
				Func: func() *Thunk {
					return ThunkFromValue(arrFrom(i + 1))
				},
			},
		}
	}
	return arrFrom(0)
}

type ArrValue struct {
	Head Value
	Tail *Thunk
}

func (v *ArrValue) GetTail() (*ArrValue, error) {
	res := v.Tail.Eval()
	if res.Error != nil {
		return nil, res.Error
	}
	if res.Value == nil {
		return nil, nil
	}
	return res.Value.(*ArrValue), nil
}

func (v *ArrValue) String() (string, error) {
	var buffer bytes.Buffer
	buffer.WriteString("[")
	if v != nil {
		head, err := v.Head.String()
		if err != nil {
			return "", err
		}
		buffer.WriteString(head)
		v, err = v.GetTail()
		if err != nil {
			return "", err
		}
		for v != nil {
			buffer.WriteString(", ")
			head, err = v.Head.String()
			if err != nil {
				return "", err
			}
			buffer.WriteString(head)
			v, err = v.GetTail()
			if err != nil {
				return "", err
			}
		}
	}
	buffer.WriteString("]")
	return buffer.String(), nil
}

func (v *ArrValue) Out() (string, error) {
	return v.String()
}

func (v *ArrValue) Inhabits(t types.Type, stack *BindingStack) (bool, error) {
	switch t := t.(type) {
	case *types.NearrType:
		if v == nil {
			return false, nil
		}
		ok, err := v.Head.Inhabits(t.HeadType, stack)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
		tail, err := v.GetTail()
		if err != nil {
			return false, err
		}
		return tail.Inhabits(t.TailType, stack)
	case *types.ArrType:
		if (types.AnyType{}).Subsumes(t.ElType) {
			return true, nil
		}
		for v != nil {
			ok, err := v.Head.Inhabits(t.ElType, stack)
			if err != nil {
				return false, err
			}
			if !ok {
				return false, nil
			}
			v, err = v.GetTail()
			if err != nil {
				return false, err
			}
		}
		return true, nil
	case types.UnionType:
		return inhabits(v, t, stack)
	case types.AnyType:
		return true, nil
	case types.TypeVariable:
		return stack.Inhabits(v, t)
	default:
		return false, nil
	}
}

func (v *ArrValue) Equal(w Value) (bool, error) {
	switch w := w.(type) {
	case *ArrValue:
		if v == nil {
			return w == nil, nil
		}
		if w == nil {
			return false, nil
		}
		ok, err := v.Head.Equal(w.Head)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
		vTail, err := v.GetTail()
		if err != nil {
			return false, err
		}
		wTail, err := w.GetTail()
		if err != nil {
			return false, err
		}
		return vTail.Equal(wTail)
	default:
		return false, nil
	}
}
