package values

import (
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

type BooleanValue struct {
	Value bool
}

func (v *BooleanValue) String() string {
	return strconv.FormatBool(v.Value)
}

///////////////////////////////////////////////////////////////////////////////

type NumberValue struct {
	Value float64
}

func (v *NumberValue) String() string {
	return strconv.FormatFloat(v.Value, 'f', -1, 64)
}

///////////////////////////////////////////////////////////////////////////////

type StringValue struct {
	Value string
}

func (v *StringValue) String() string {
	return fmt.Sprintf("%q", v.Value)
}

///////////////////////////////////////////////////////////////////////////////
