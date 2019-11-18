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

func (v *ArrValue) String() (string, error) {
	var buffer bytes.Buffer
	buffer.WriteString("[")
	err := v.Eval()
	if err != nil {
		return "", err
	}
	if v.Head != nil {
		head, err := v.Head.String()
		if err != nil {
			return "", err
		}
		buffer.WriteString(head)
		v = v.Tail
		err = v.Eval()
		if err != nil {
			return "", err
		}
		for v.Head != nil {
			buffer.WriteString(", ")
			head, err := v.Head.String()
			if err != nil {
				return "", err
			}
			buffer.WriteString(head)
			v = v.Tail
			err = v.Eval()
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
		err := v.Eval()
		if err != nil {
			return false, err
		}
		if v.Head == nil {
			return false, nil
		}
		if ok, err := v.Head.Inhabits(t.HeadType, stack); !ok {
			return false, err
		}
		return v.Tail.Inhabits(t.TailType, stack)
	case *types.ArrType:
		if (types.AnyType{}).Subsumes(t.ElType) {
			return true, nil
		}
		err := v.Eval()
		if err != nil {
			return false, err
		}
		for v.Head != nil {
			if ok, err := v.Head.Inhabits(t.ElType, stack); !ok {
				return false, err
			}
			v = v.Tail
			err := v.Eval()
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
		err := v.Eval()
		if err != nil {
			return false, err
		}
		err = w.Eval()
		if err != nil {
			return false, err
		}
		if v.Head == nil {
			return w.Head == nil, nil
		}
		if w.Head == nil {
			return false, nil
		}
		if ok, err := v.Head.Equal(w.Head); !ok {
			return false, err
		}
		return v.Tail.Equal(w.Tail)
	default:
		return false, nil
	}
}
