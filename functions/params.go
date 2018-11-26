package functions

import (
	"fmt"
	"strings"

	"github.com/texttheater/bach/actions"
	"github.com/texttheater/bach/types"
)

type Param struct {
	InputType   types.Type
	Name        string
	Params      []*Param
	OutputType  types.Type
	ActionStack *ActionStack
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

func (p *Param) String() string {
	paramStrings := make([]string, 0, len(p.Params))
	for _, param := range p.Params {
		paramStrings = append(paramStrings, param.String())
	}
	paramsString := strings.Join(paramStrings, ",")
	return fmt.Sprintf("for %s %s(%s) %s", p.InputType, p.Name, paramsString, p.OutputType)
}

type ActionStack struct {
	Head actions.Action
	Tail *ActionStack
}

func (s *ActionStack) Push(action actions.Action) *ActionStack {
	return &ActionStack{
		Head: action,
		Tail: s,
	}
}
