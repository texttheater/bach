package parameters

import (
	"bytes"

	"github.com/texttheater/bach/types"
)

type Parameter struct {
	InputType  types.Type
	Name       string
	Params     []*Parameter
	OutputType types.Type
}

func (this *Parameter) Subsumes(that *Parameter) bool {
	if len(this.Params) != len(that.Params) {
		return false
	}
	if !that.InputType.Subsumes(this.InputType) {
		return false
	}
	if !this.OutputType.Subsumes(that.OutputType) {
		return false
	}
	for i, thatParam := range that.Params {
		if !thatParam.Subsumes(this.Params[i]) {
			return false
		}
	}
	return true
}

func (p *Parameter) String() string {
	buffer := bytes.Buffer{}
	if !p.InputType.Subsumes(types.AnyType) || len(p.Params) > 0 {
		buffer.WriteString("for ")
		buffer.WriteString(p.InputType.String())
		buffer.WriteString(" ")
	}
	buffer.WriteString(p.Name)
	if len(p.Params) > 0 {
		buffer.WriteString("(")
		buffer.WriteString(p.Params[0].String())
		for _, param := range p.Params[1:] {
			buffer.WriteString(",")
			buffer.WriteString(param.String())
		}
		buffer.WriteString(")")
	}
	buffer.WriteString(" ")
	buffer.WriteString(p.OutputType.String())
	return buffer.String()
}
