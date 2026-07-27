package states

import (
	"bytes"

	"github.com/texttheater/bach/types"
)

func ArrFromSlice(elements []Value) *ArrValue {
	var arrFrom func(i int) *ArrValue
	arrFrom = func(i int) *ArrValue {
		if i == len(elements) {
			return nil
		}
		return &ArrValue{
			Head: ThunkFromValue(elements[i]),
			Tail: ThunkFromFunc(func() (Value, error) {
				return arrFrom(i + 1), nil
			}),
		}
	}
	return arrFrom(0)
}

type ArrValue struct {
	Head *Thunk
	Tail *Thunk
}

func (v *ArrValue) Repr() (string, error) {
	buffer := bytes.Buffer{}
	buffer.WriteString("[")
	if v != nil {
		head, err := v.Head.Eval()
		if err != nil {
			return "", err
		}
		headRepr, err := head.Repr()
		if err != nil {
			return "", err
		}
		buffer.WriteString(headRepr)
		v, err = v.Tail.EvalArr()
		if err != nil {
			return "", err
		}
		for v != nil {
			buffer.WriteString(", ")
			head, err = v.Head.Eval()
			if err != nil {
				return "", err
			}
			headRepr, err := head.Repr()
			if err != nil {
				return "", err
			}
			buffer.WriteString(headRepr)
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
		head, err := v.Head.Eval()
		if err != nil {
			return nil, err
		}
		data, err := head.Data()
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
		head, err := v.Head.Eval()
		if err != nil {
			return false, err
		}
		ok, err := head.Inhabits(t.Head, stack)
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
			head, err := v.Head.Eval()
			if err != nil {
				return false, err
			}
			ok, err := head.Inhabits(t.El, stack)
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
		vHead, err := v.Head.Eval()
		if err != nil {
			return false, err
		}
		wHead, err := w.Head.Eval()
		if err != nil {
			return false, err
		}
		ok, err := vHead.Equal(wHead)
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

func IterFromThunk(t *Thunk) func() (*Thunk, bool, error) {
	return func() (*Thunk, bool, error) {
		arr, err := t.EvalArr()
		if err != nil {
			return nil, false, err
		}
		if arr == nil {
			return nil, false, nil
		}
		t = arr.Tail
		return arr.Head, true, nil
	}
}

func ThunkFromIter(iter func() (*Thunk, bool, error)) *Thunk {
	head, ok, err := iter()
	if err != nil {
		return ThunkFromError(err)
	}
	if !ok {
		return ThunkFromValue((*ArrValue)(nil))
	}
	return ThunkFromValue(&ArrValue{
		Head: head,
		Tail: ThunkFromIter(iter),
	})
}

func SliceFromThunk(t *Thunk) ([]*Thunk, error) {
	var slice []*Thunk
	iter := IterFromThunk(t)
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

func NumsFromThunk(t *Thunk) ([]float64, error) {
	var slice []float64
	iter := IterFromThunk(t)
	for {
		el, ok, err := iter()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		num, err := el.EvalNum()
		if err != nil {
			return nil, err
		}
		slice = append(slice, num)
	}
	return slice, nil
}

func ThunkFromSlice(slice []*Thunk) *Thunk {
	i := 0
	iter := func() (*Thunk, bool, error) {
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
	iter := func() (*Thunk, bool, error) {
		thk := <-c
		val, err := thk.Eval()
		if err != nil {
			return nil, false, err
		}
		if val == nil {
			return nil, false, nil
		}
		return thk, true, nil
	}
	return ThunkFromIter(iter)
}
