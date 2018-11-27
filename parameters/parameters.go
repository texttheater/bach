package parameters

import (
	"fmt"
	"strings"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type Parameter struct {
	InputType   types.Type
	Name        string
	Params      []*Parameter
	OutputType  types.Type
	ActionStack *ActionStack
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
	paramStrings := make([]string, 0, len(p.Params))
	for _, param := range p.Params {
		paramStrings = append(paramStrings, param.String())
	}
	paramsString := strings.Join(paramStrings, ",")
	return fmt.Sprintf("for %s %s(%s) %s", p.InputType, p.Name, paramsString, p.OutputType)
}

type ActionStack struct {
	Head states.Action
	Tail *ActionStack
}

func (s *ActionStack) Push(action states.Action) *ActionStack {
	return &ActionStack{
		Head: action,
		Tail: s,
	}
}
