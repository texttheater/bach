package parameters

import (
	"bytes"

	"github.com/texttheater/bach/types"
)

type Parameter struct {
	InputType  types.Type
	Params     []*Parameter
	OutputType types.Type
}

func SimpleParam(outputType types.Type) *Parameter {
	return &Parameter{
		InputType:  types.Any{},
		Params:     nil,
		OutputType: outputType,
	}
}

func (p *Parameter) Subsumes(q *Parameter) bool {
	if len(p.Params) != len(q.Params) {
		return false
	}
	if !q.InputType.Subsumes(p.InputType) {
		return false
	}
	if !p.OutputType.Subsumes(q.OutputType) {
		return false
	}
	for i, otherParam := range q.Params {
		if !otherParam.Subsumes(p.Params[i]) {
			return false
		}
	}
	return true
}

func (p *Parameter) Equivalent(q *Parameter) bool {
	return p.Subsumes(q) && q.Subsumes(p)
}

func (p *Parameter) Instantiate(bindings map[string]types.Type) *Parameter {
	inputType := p.InputType.Instantiate(bindings)
	var params []*Parameter
	if p.Params != nil {
		params = make([]*Parameter, len(p.Params))
		for i, param := range p.Params {
			params[i] = param.Instantiate(bindings)
		}
	}
	outputType := p.OutputType.Instantiate(bindings)
	return &Parameter{
		InputType:  inputType,
		Params:     params,
		OutputType: outputType,
	}
}

func (p Parameter) String() string {
	buffer := bytes.Buffer{}
	if !p.InputType.Subsumes(types.Any{}) || len(p.Params) > 0 {
		buffer.WriteString("for ")
		buffer.WriteString(p.InputType.String())
		buffer.WriteString(" ")
	}
	if len(p.Params) > 0 {
		buffer.WriteString("(")
		buffer.WriteString(p.Params[0].String())
		for _, param := range p.Params[1:] {
			buffer.WriteString(", ")
			buffer.WriteString(param.String())
		}
		buffer.WriteString(") ")
	}
	buffer.WriteString(p.OutputType.String())
	return buffer.String()
}
