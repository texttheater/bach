package expressions

import (
	"regexp"
	"strconv"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/shapes"
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

func (x RegexpExpression) Typecheck(inputShape shapes.Shape, params []*params.Param) (shapes.Shape, states.Action, *states.IDStack, error) {
	if len(params) > 0 {
		return shapes.Shape{}, nil, nil, errors.TypeError(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	if !(types.Str{}).Subsumes(inputShape.Type) {
		return shapes.Shape{}, nil, nil, errors.TypeError(
			errors.Code(errors.RegexpWantsString),
			errors.Pos(x.Pos),
			errors.WantType(types.Str{}),
			errors.GotType(inputShape.Type),
		)
	}
	submatchType := types.NewUnion(types.Null{}, types.Str{})
	propTypeMap := make(map[string]types.Type)
	propTypeMap["0"] = types.Str{}
	propTypeMap["start"] = types.Num{}
	subexpNames := x.Regexp.SubexpNames()
	for i := 1; i < len(subexpNames); i++ {
		name := subexpNames[i]
		propTypeMap[strconv.Itoa(i)] = submatchType
		if name != "" {
			propTypeMap[name] = submatchType
		}
	}
	matchType := types.Obj{
		Props: propTypeMap,
		Rest:  types.Void{},
	}
	outputShape := shapes.Shape{
		Type:  types.NewUnion(types.Null{}, matchType),
		Stack: inputShape.Stack,
	}
	action := func(inputState states.State, args []states.Action) states.State {
		inputString, err := inputState.Thunk.EvalStr()
		if err != nil {
			return inputState.Replace(states.ThunkFromError(err))
		}
		match := x.Regexp.FindStringSubmatchIndex(inputString)
		if match == nil {
			return inputState.Replace(states.ThunkFromValue(states.NullValue{}))
		}
		propThunkMap := make(map[string]*states.Thunk)
		if propTypeMap["start"].Subsumes(types.Num{}) {
			propThunkMap["start"] = states.ThunkFromValue(states.NumValue(match[0]))
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
			propThunkMap[strconv.Itoa(i)] = states.ThunkFromValue(submatch)
			if name != "" {
				propThunkMap[name] = states.ThunkFromValue(submatch)
			}
		}
		return inputState.Replace(states.ThunkFromValue(states.ObjValue(propThunkMap)))
	}
	return outputShape, action, nil, nil
}
