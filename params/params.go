package params

import (
	"bytes"

	"github.com/texttheater/bach/types"
)

type Param struct {
	InputType   types.Type
	Name        string
	Description string
	Params      []*Param
	OutputType  types.Type
}

func SimpleParam(Name string, Description string, outputType types.Type) *Param {
	return &Param{
		InputType:   types.Any{},
		Name:        Name,
		Description: Description,
		Params:      nil,
		OutputType:  outputType,
	}
}

func (p *Param) Subsumes(q *Param) bool {
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

func (p *Param) Equivalent(q *Param) bool {
	return p.Subsumes(q) && q.Subsumes(p)
}

func (p *Param) Instantiate(bindings map[string]types.Type) *Param {
	inputType := p.InputType.Instantiate(bindings)
	var params []*Param
	if p.Params != nil {
		params = make([]*Param, len(p.Params))
		for i, param := range p.Params {
			params[i] = param.Instantiate(bindings)
		}
	}
	outputType := p.OutputType.Instantiate(bindings)
	return &Param{
		InputType:  inputType,
		Name:       p.Name,
		Params:     params,
		OutputType: outputType,
	}
}

func (p Param) String() string {
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
