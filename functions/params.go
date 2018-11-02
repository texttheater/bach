package functions

import (
	"github.com/texttheater/bach/types"
)

type Param struct {
	InputType  types.Type
	Name       string
	Params     []*Param
	OutputType types.Type
}

func (this *Param) Subsumes(that *Param) bool {
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

func (p *Param) DummyFunction() Function {
	return Function{
		InputType: p.InputType,
		Name: p.Name,
		Params: p.Params,
		OutputType: p.OutputType,
		Action: &Action{
			Execute: nil,
		},
	}
}
