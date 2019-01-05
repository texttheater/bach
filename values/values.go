package values

import (
	"bytes"
	"fmt"
	"strconv"
)

///////////////////////////////////////////////////////////////////////////////

type Value interface {
	String() string
	Iter() <-chan Value
}

///////////////////////////////////////////////////////////////////////////////

type NullValue struct {
}

func (v *NullValue) String() string {
	return "null"
}

func (v *NullValue) Iter() <-chan Value {
	panic(fmt.Sprintf("%s is not a sequence", v))
}

///////////////////////////////////////////////////////////////////////////////

type BoolValue struct {
	Value bool
}

func (v *BoolValue) String() string {
	return strconv.FormatBool(v.Value)
}

func (v *BoolValue) Iter() <-chan Value {
	panic(fmt.Sprintf("%s is not a sequence", v))
}

///////////////////////////////////////////////////////////////////////////////

type NumValue struct {
	Value float64
}

func (v *NumValue) String() string {
	return strconv.FormatFloat(v.Value, 'f', -1, 64)
}

func (v *NumValue) Iter() <-chan Value {
	panic(fmt.Sprintf("%s is not a sequence", v))
}

///////////////////////////////////////////////////////////////////////////////

type StrValue struct {
	Value string
}

func (v *StrValue) String() string {
	return fmt.Sprintf("%q", v.Value)
}

func (v *StrValue) Iter() <-chan Value {
	panic(fmt.Sprintf("%s is not a sequence", v))
}

///////////////////////////////////////////////////////////////////////////////

type ArrValue struct {
	ElementValues []Value
}

func (v *ArrValue) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("[")
	if len(v.ElementValues) > 0 {
		buffer.WriteString(v.ElementValues[0].String())
		for _, elValue := range v.ElementValues[1:] {
			buffer.WriteString(", ")
			buffer.WriteString(elValue.String())
		}
	}
	buffer.WriteString("]")
	return buffer.String()
}

func (v *ArrValue) Iter() <-chan Value {
	channel := make(chan Value)
	go func() {
		for _, el := range v.ElementValues {
			channel <- el
		}
		close(channel)
	}()
	return channel
}

///////////////////////////////////////////////////////////////////////////////
