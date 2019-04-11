package values

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/texttheater/bach/types"
)

type Value interface {
	Type() types.Type
	String() string
	Out() string
	Iter() <-chan Value
	Arr() []Value
	Bool() bool
	Num() float64
	Obj() map[string]Value
	Reader() io.Reader
	Str() string
}

func ArrValue(arr []Value) Value {
	t := types.VoidArrType
	for i := len(arr) - 1; i >= 0; i-- {
		t = types.NearrType(arr[i].Type(), t)
	}
	return value{t: t, arr: arr}
}

func BoolValue(b bool) Value {
	return value{t: types.BoolType(), b: b}
}

var nullValue Value = value{t: types.NullType()}

func NullValue() Value {
	return nullValue
}

func NumValue(num float64) Value {
	return value{t: types.NumType(), num: num}
}

func ObjValue(obj map[string]Value) Value {
	propTypeMap := make(map[string]types.Type)
	for k, v := range obj {
		propTypeMap[k] = v.Type()
	}
	t := types.ObjType(propTypeMap)
	return value{t: t, obj: obj}
}

func ReaderValue(reader io.Reader) Value {
	return value{t: types.ReaderType(), reader: reader}
}

func SeqValue(elementType types.Type, seq chan Value) Value {
	return value{t: types.SeqType(elementType), seq: seq}
}

func StrValue(str string) Value {
	return value{t: types.StrType, str: str}
}

type value struct {
	t      types.Type
	arr    []Value
	b      bool
	num    float64
	obj    map[string]Value
	reader io.Reader
	seq    chan Value
	str    string
}

func (v value) Type() types.Type {
	return v.t
}

func (v value) String() string {
	if types.AnyArrType.Subsumes(v.t) {
		var buffer bytes.Buffer
		buffer.WriteString("[")
		if len(v.arr) > 0 {
			buffer.WriteString(v.arr[0].String())
			for _, elValue := range v.arr[1:] {
				buffer.WriteString(", ")
				buffer.WriteString(elValue.String())
			}
		}
		buffer.WriteString("]")
		return buffer.String()
	}
	if types.BoolType().Subsumes(v.t) {
		return strconv.FormatBool(v.b)
	}
	if types.NumType().Subsumes(v.t) {
		return strconv.FormatFloat(v.num, 'f', -1, 64)
	}
	if types.AnyObjType.Subsumes(v.t) {
		var buffer bytes.Buffer
		buffer.WriteString("{")
		firstWritten := false
		for k, w := range v.obj {
			if firstWritten {
				buffer.WriteString(", ")
			}
			buffer.WriteString(k)
			buffer.WriteString(": ")
			buffer.WriteString(w.String())
			firstWritten = true
		}
		buffer.WriteString("}")
		return buffer.String()
	}
	if types.ReaderType().Subsumes(v.t) {
		return "<reader>"
	}
	if types.AnySeqType.Subsumes(v.t) {
		return "<seq>"
	}
	if types.StrType.Subsumes(v.t) {
		return fmt.Sprintf("%q", v.str)
	}
	panic("invalid value")
}

func (v value) Out() string {
	if types.AnyArrType.Subsumes(v.t) {
		var buffer bytes.Buffer
		buffer.WriteString("[")
		if len(v.arr) > 0 {
			buffer.WriteString(v.arr[0].String())
			for _, elValue := range v.arr[1:] {
				buffer.WriteString(", ")
				buffer.WriteString(elValue.String())
			}
		}
		buffer.WriteString("]")
		return buffer.String()
	}
	if types.BoolType().Subsumes(v.t) {
		return strconv.FormatBool(v.b)
	}
	if types.NumType().Subsumes(v.t) {
		return strconv.FormatFloat(v.num, 'f', -1, 64)
	}
	if types.AnyObjType.Subsumes(v.t) {
		var buffer bytes.Buffer
		buffer.WriteString("{")
		firstWritten := false
		for k, w := range v.obj {
			if firstWritten {
				buffer.WriteString(", ")
			}
			buffer.WriteString(k)
			buffer.WriteString(": ")
			buffer.WriteString(w.String())
			firstWritten = true
		}
		buffer.WriteString("}")
		return buffer.String()
	}
	if types.ReaderType().Subsumes(v.t) {
		return "<reader>"
	}
	if types.AnySeqType.Subsumes(v.t) {
		return "<seq>"
	}
	if types.StrType.Subsumes(v.t) {
		return v.str
	}
	panic("invalid value")
}

func (v value) Iter() <-chan Value {
	if types.AnyArrType.Subsumes(v.t) {
		channel := make(chan Value)
		go func() {
			for _, el := range v.arr {
				channel <- el
			}
			close(channel)
		}()
		return channel
	}
	if types.AnySeqType.Subsumes(v.t) {
		// TODO safeguard against iterating twice?
		return v.seq
	}
	panic(fmt.Sprintf("%s is not a sequence", v))
}

func (v value) Arr() []Value {
	return v.arr
}

func (v value) Bool() bool {
	return v.b
}

func (v value) Num() float64 {
	return v.num
}

func (v value) Obj() map[string]Value {
	return v.obj
}

func (v value) Reader() io.Reader {
	return v.reader
}

func (v value) Str() string {
	return v.str
}
