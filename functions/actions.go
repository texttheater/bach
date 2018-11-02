package functions

import (
	"github.com/texttheater/bach/values"
)

type Action struct {
	Execute func(inputValue values.Value, args []*Action) values.Value
}

func (a *Action) SetArg(arg *Action) *Action {
	return &Action{
		Execute: func(inputValue values.Value, outerArgs []*Action) values.Value {
			args := make([]*Action, 0, len(outerArgs)+1)
			args = append(args, arg)
			for _, outerArg := range outerArgs {
				args = append(args, outerArg)
			}
			return a.Execute(inputValue, args)
		},
	}
}
