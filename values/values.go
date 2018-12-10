package values

import (
	"bytes"
	"fmt"
	"strconv"
)

///////////////////////////////////////////////////////////////////////////////

type Value interface {
	String() string
}

///////////////////////////////////////////////////////////////////////////////

type NullValue struct {
}

func (v *NullValue) String() string {
	return "null"
}

///////////////////////////////////////////////////////////////////////////////

type BoolValue struct {
	Value bool
}

func (v *BoolValue) String() string {
	return strconv.FormatBool(v.Value)
}

///////////////////////////////////////////////////////////////////////////////

type NumValue struct {
	Value float64
}

func (v *NumValue) String() string {
	return strconv.FormatFloat(v.Value, 'f', -1, 64)
}

///////////////////////////////////////////////////////////////////////////////

type StrValue struct {
	Value string
}

func (v *StrValue) String() string {
	return fmt.Sprintf("%q", v.Value)
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

///////////////////////////////////////////////////////////////////////////////
