package functions

import (
	"github.com/texttheater/bach/types"
)

type Param struct {
	InputType  types.Type
	Name       string // TODO do we actually use this?
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
