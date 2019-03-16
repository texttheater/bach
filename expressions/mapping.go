package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type MappingExpression struct {
	Pos  lexer.Position
	Body Expression
}

func (x *MappingExpression) Typecheck(inputShape functions.Shape, params []*functions.Parameter) (functions.Shape, functions.Action, error) {
	// make sure we got no parameters
	if len(params) > 0 {
		return zeroShape, nil, errors.E(
			errors.Kind(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	// make sure the input type is a sequence type
	if !types.AnySeqType.Subsumes(inputShape.Type) {
		return zeroShape, nil, errors.E(
			errors.Kind(errors.MappingRequiresSeqType),
			errors.Pos(x.Pos),
			errors.WantType(types.AnySeqType),
			errors.GotType(inputShape.Type),
		)
	}
	// typecheck body
	bodyInputShape := functions.Shape{
		Type:        inputShape.Type.ElementType(),
		FuncerStack: inputShape.FuncerStack,
	}
	bodyOutputShape, bodyAction, err := x.Body.Typecheck(bodyInputShape, nil)
	if err != nil {
		return zeroShape, nil, err
	}
	// create output shape
	outputShape := functions.Shape{
		Type:        types.SeqType(bodyOutputShape.Type),
		FuncerStack: inputShape.FuncerStack,
	}
	// create action
	action := func(inputState functions.State, args []functions.Action) functions.State {
		channel := make(chan values.Value)
		go func() {
			for el := range inputState.Value.Iter() {
				bodyInputState := functions.State{
					Value: el,
					Stack: inputState.Stack,
				}
				bodyOutputState := bodyAction(bodyInputState, nil)
				channel <- bodyOutputState.Value
			}
			close(channel)
		}()
		return functions.State{
			Value: values.SeqValue(channel),
			Stack: inputState.Stack,
		}
	}
	return outputShape, action, nil
}
