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

type RegexpFindFirstExpression struct {
	Pos    lexer.Position
	Regexp *regexp.Regexp
}

func (x RegexpFindFirstExpression) Position() lexer.Position {
	return x.Pos
}

func (x RegexpFindFirstExpression) Typecheck(inputShape Shape, params []*Parameter) (Shape, states.Action, error) {
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
	propTypeMap["start"] = types.NumType{}
	for i, name := range x.Regexp.SubexpNames() {
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
	action := func(inputState states.State, args []states.Action) states.State {
		inputString := string(inputState.Value.(values.StrValue))
		match := x.Regexp.FindStringSubmatchIndex(inputString)
		if match == nil {
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
			fromIndex := match[2*i]
			toIndex := match[2*i+1]
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

type RegexpFindAllExpression struct {
	Pos    lexer.Position
	Regexp *regexp.Regexp
}

func (x RegexpFindAllExpression) Position() lexer.Position {
	return x.Pos
}

func (x RegexpFindAllExpression) Typecheck(inputShape Shape, params []*Parameter) (Shape, states.Action, error) {
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
	propTypeMap["start"] = types.NumType{}
	for i, name := range x.Regexp.SubexpNames() {
		propTypeMap[strconv.Itoa(i)] = submatchType
		if name != "" {
			propTypeMap[name] = submatchType
		}
	}
	matchType := types.NewObjType(propTypeMap)
	outputShape := Shape{
		Type:  &types.SeqType{matchType},
		Stack: inputShape.Stack,
	}
	action := func(inputState states.State, args []states.Action) states.State {
		inputString := string(inputState.Value.(values.StrValue))
		seq := values.SeqValue{
			ElementType: matchType,
			Channel:     make(chan values.Value),
		}
		go func() {
			for _, match := range x.Regexp.FindAllStringSubmatchIndex(inputString, -1) {
				propValueMap := make(map[string]values.Value)
				if propTypeMap["start"].Subsumes(types.NumType{}) {
					propValueMap["start"] = values.NumValue(match[0])
				}
				for i, name := range x.Regexp.SubexpNames() {
					fromIndex := match[2*i]
					toIndex := match[2*i+1]
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
				seq.Channel <- values.ObjValue(propValueMap)
			}
			close(seq.Channel)
		}()
		return states.State{
			Value: seq,
			Stack: inputState.Stack,
		}
	}
	return outputShape, action, nil
}
