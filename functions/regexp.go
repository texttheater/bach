package functions

import (
	"regexp"
	"strconv"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type RegexpMatchExpression struct {
	Pos    lexer.Position
	Regexp *regexp.Regexp
}

func (x RegexpMatchExpression) Position() lexer.Position {
	return x.Pos
}

func (x RegexpMatchExpression) Typecheck(inputShape Shape, params []*Parameter) (Shape, states.Action, error) {
	if !(types.StrType{}).Subsumes(inputShape.Type) {
		return Shape{}, nil, errors.E(
			errors.Code(errors.RegexpWantsString),
			errors.Pos(x.Pos),
			errors.WantType(types.StrType{}),
			errors.GotType(inputShape.Type),
		)
	}
	submatchType := types.Union(types.NullType{}, types.StrType{})
	propTypeMap := make(map[string]types.Type)
	propTypeMap["start"] = types.NumType{}
	for i, name := range x.Regexp.SubexpNames() {
		propTypeMap[strconv.Itoa(i)] = submatchType
		if name != "" {
			propTypeMap[name] = submatchType
		}
	}
	outputShape := Shape{
		Type:  types.Union(types.NullType{}, types.NewObjType(propTypeMap)),
		Stack: inputShape.Stack,
	}
	action := func(inputState states.State, args []states.Action) states.State {
		inputString := string(inputState.Value.(values.StrValue))
		indexes := x.Regexp.FindStringSubmatchIndex(inputString)
		if indexes == nil {
			return states.State{
				Value: values.NullValue{},
				Stack: inputState.Stack,
			}
		}
		propValueMap := make(map[string]values.Value)
		if propTypeMap["start"].Subsumes(types.NumType{}) {
			propValueMap["start"] = values.NumValue(match[0])
		}
		for i, name := range x.Regexp.SubexpNames() {
			fromIndex := indexes[2*i]
			toIndex := indexes[2*i+1]
			var submatch values.Value
			if fromIndex == -1 {
				submatch = values.NullValue{}
			} else {
				submatch = values.StrValue(inputString[fromIndex:toIndex])
			}
			propValueMap[strconv.Itoa(i)] = submatch
			if name != "" {
				propValueMap[name] = submatch
			}
		}
		return states.State{
			Value: values.ObjValue(propValueMap),
			Stack: inputState.Stack,
		}
	}
	return outputShape, action, nil
}
