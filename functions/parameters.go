package functions

import (
	"fmt"
	"strings"

	"github.com/texttheater/bach/types"
)

type Parameter struct {
	InputType  types.Type
	Parameters []*Parameter
	OutputType types.Type
}

func (p *Parameter) String() string {
	pars := make([]string, 0, len(p.Parameters))
	for _, par := range p.Parameters {
		pars = append(pars, par.String())
	}
	return fmt.Sprintf("for %s (%s) %s", p.InputType,
		strings.Join(pars, ""), p.OutputType)
}

func (p *Parameter) Subsumes(q *Parameter) bool {
	if len(p.Parameters) != len(q.Parameters) {
		return false
	}
	if !q.InputType.Subsumes(p.InputType) {
		return false
	}
	for i, pPar := range p.Parameters {
		if !q.Parameters[i].Subsumes(pPar) {
			return false
		}
	}
	if !p.OutputType.Subsumes(q.OutputType) {
		return false
	}
	return true
}

func SimpleParameter(OutputType types.Type) *Parameter {
	return &Parameter{
		&types.AnyType{},
		[]*Parameter{},
		OutputType,
	}
}
