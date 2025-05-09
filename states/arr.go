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

func (v *ArrValue) Repr() (string, error) {
	buffer := bytes.Buffer{}
	buffer.WriteString("[")
	if v != nil {
		head, err := v.Head.Repr()
		if err != nil {
			return "", err
		}
		buffer.WriteString(head)
		v, err = v.Tail.EvalArr()
		if err != nil {
			return "", err
		}
		for v != nil {
			buffer.WriteString(", ")
			head, err = v.Head.Repr()
			if err != nil {
				return "", err
			}
			buffer.WriteString(head)
			v, err = v.Tail.EvalArr()
			if err != nil {
				return "", err
			}
		}
	}
	buffer.WriteString("]")
	return buffer.String(), nil
}

func (v *ArrValue) Str() (string, error) {
	return v.Repr()
}

func (v *ArrValue) Data() (any, error) {
	res := make([]any, 0)
	for v != nil {
		data, err := v.Head.Data()
		if err != nil {
			return nil, err
		}
		res = append(res, data)
		v, err = v.Tail.EvalArr()
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (v *ArrValue) Inhabits(t types.Type, stack *BindingStack) (bool, error) {
	switch t := t.(type) {
	case *types.Nearr:
		if v == nil {
			return false, nil
		}
		ok, err := v.Head.Inhabits(t.Head, stack)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
		tail, err := v.Tail.EvalArr()
		if err != nil {
			return false, err
		}
		return tail.Inhabits(t.Tail, stack)
	case *types.Arr:
		if (types.Any{}).Subsumes(t.El) {
			return true, nil
		}
		for v != nil {
			ok, err := v.Head.Inhabits(t.El, stack)
			if err != nil {
				return false, err
			}
			if !ok {
				return false, nil
			}
			v, err = v.Tail.EvalArr()
			if err != nil {
				return false, err
			}
		}
		return true, nil
	case types.Union:
		return inhabits(v, t, stack)
	case types.Any:
		return true, nil
	case types.Var:
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
		vTail, err := v.Tail.EvalArr()
		if err != nil {
			return false, err
		}
		wTail, err := w.Tail.EvalArr()
		if err != nil {
			return false, err
		}
		return vTail.Equal(wTail)
	default:
		return false, nil
	}
}

func IterFromValue(v Value) func() (Value, bool, error) {
	return func() (Value, bool, error) {
		arr := v.(*ArrValue)
		if arr == nil {
			return nil, false, nil
		}
		var err error
		v, err = arr.Tail.EvalArr()
		if err != nil {
			return nil, false, err
		}
		return arr.Head, true, nil
	}
}

func ThunkFromIter(iter func() (Value, bool, error)) *Thunk {
	value, ok, err := iter()
	if err != nil {
		return ThunkFromError(err)
	}
	if !ok {
		return ThunkFromValue((*ArrValue)(nil))
	}
	return ThunkFromValue(&ArrValue{
		Head: value,
		Tail: &Thunk{
			Func: func() *Thunk {
				return ThunkFromIter(iter)
			},
		},
	})
}

func SliceFromValue(v Value) ([]Value, error) {
	var slice []Value
	iter := IterFromValue(v)
	for {
		el, ok, err := iter()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		slice = append(slice, el)
	}
	return slice, nil
}

func NumsFromValue(v Value) ([]float64, error) {
	var slice []float64
	iter := IterFromValue(v)
	for {
		el, ok, err := iter()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		slice = append(slice, float64(el.(NumValue)))
	}
	return slice, nil
}

func ThunkFromSlice(slice []Value) *Thunk {
	i := 0
	iter := func() (Value, bool, error) {
		if i < len(slice) {
			el := slice[i]
			i++
			return el, true, nil
		}
		return nil, false, nil
	}
	return ThunkFromIter(iter)
}

func ThunkFromChannel(c <-chan *Thunk) *Thunk {
	iter := func() (Value, bool, error) {
		thk := <-c
		val, err := thk.Eval()
		if err != nil {
			return nil, false, err
		}
		if val == nil {
			return nil, false, nil
		}
		return val, true, nil
	}
	return ThunkFromIter(iter)
}
