package functions

import (
	"regexp"
	"strconv"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type RegexpExpression struct {
	Pos    lexer.Position
	Regexp *regexp.Regexp
}

func (x RegexpExpression) Position() lexer.Position {
	return x.Pos
}

func (x RegexpExpression) Typecheck(inputShape Shape, params []*Parameter) (Shape, states.Action, error) {
	if len(params) > 0 {
		return Shape{}, nil, errors.E(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
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
	propTypeMap["0"] = types.StrType{}
	propTypeMap["start"] = types.NumType{}
	subexpNames := x.Regexp.SubexpNames()
	for i := 1; i < len(subexpNames); i++ {
		name := subexpNames[i]
		propTypeMap[strconv.Itoa(i)] = submatchType
		if name != "" {
			propTypeMap[name] = submatchType
		}
	}
	matchType := types.NewObjType(propTypeMap)
	outputShape := Shape{
		Type:  types.Union(types.NullType{}, matchType),
		Stack: inputShape.Stack,
	}
	action := func(inputState states.State, args []states.Action) *states.Thunk {
		inputString := string(inputState.Value.(states.StrValue))
		match := x.Regexp.FindStringSubmatchIndex(inputString)
		if match == nil {
			return &states.Thunk{State: states.State{
				Value:     states.NullValue{},
				Stack:     inputState.Stack,
				TypeStack: inputState.TypeStack,
			}}

		}
		propValueMap := make(map[string]states.Value)
		if propTypeMap["start"].Subsumes(types.NumType{}) {
			propValueMap["start"] = states.NumValue(match[0])
		}
		for i, name := range x.Regexp.SubexpNames() {
			fromIndex := match[2*i]
			toIndex := match[2*i+1]
			var submatch states.Value
			if fromIndex == -1 {
				submatch = states.NullValue{}
			} else {
				submatch = states.StrValue(inputString[fromIndex:toIndex])
			}
			propValueMap[strconv.Itoa(i)] = submatch
			if name != "" {
				propValueMap[name] = submatch
			}
		}
		return &states.Thunk{State: states.State{
			Value:     states.ObjValue(propValueMap),
			Stack:     inputState.Stack,
			TypeStack: inputState.TypeStack,
		}}

	}
	return outputShape, action, nil
}
