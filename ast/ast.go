/*
Package ast implements Bach's abstract syntax trees.

An alternative name for this package would be: expressions. Because that's what
an AST is, an expression consisting of subexpressions.
*/
package ast

import (
	"fmt"
	//"os"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

///////////////////////////////////////////////////////////////////////////////

var zeroShape = functions.Shape{}

var boolType = &types.BoolType{}

var seqType = &types.SeqType{&types.AnyType{}}

///////////////////////////////////////////////////////////////////////////////

type Expression interface {
	Typecheck(inputShape functions.Shape, params []*functions.Parameter) (outputShape functions.Shape, action functions.Action, err error)
}

///////////////////////////////////////////////////////////////////////////////

type ConstantExpression struct {
	Pos   lexer.Position
	Type  types.Type
	Value values.Value
}

func (x *ConstantExpression) Typecheck(inputShape functions.Shape, params []*functions.Parameter) (functions.Shape, functions.Action, error) {
	if len(params) > 0 {
		return zeroShape, nil, errors.E(
			errors.Kind(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	outputShape := functions.Shape{x.Type, inputShape.FuncerStack}
	action := func(inputState functions.State, args []functions.Action) functions.State {
		return functions.State{
			Value: x.Value,
			Stack: inputState.Stack,
		}
	}
	return outputShape, action, nil
}

///////////////////////////////////////////////////////////////////////////////

type ArrExpression struct {
	Pos      lexer.Position
	Elements []Expression
}

func (x *ArrExpression) Typecheck(inputShape functions.Shape, params []*functions.Parameter) (functions.Shape, functions.Action, error) {
	if len(params) > 0 {
		return zeroShape, nil, errors.E(
			errors.Kind(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	var elementType types.Type = &types.AnyType{}
	elementActions := make([]functions.Action, 0, len(x.Elements))
	for i, elExpression := range x.Elements {
		elOutputShape, elAction, err := elExpression.Typecheck(inputShape, nil)
		if err != nil {
			return zeroShape, nil, err
		}
		if i == 0 {
			elementType = elOutputShape.Type // HACK
		} else {
			elementType = types.Disjoin(elementType, elOutputShape.Type)
		}
		elementActions = append(elementActions, elAction)
	}
	outputShape := functions.Shape{
		Type:        &types.ArrType{elementType},
		FuncerStack: inputShape.FuncerStack,
	}
	action := func(inputState functions.State, args []functions.Action) functions.State {
		elementValues := make([]values.Value, 0, len(elementActions))
		for _, elAction := range elementActions {
			elValue := elAction(inputState, nil).Value
			elementValues = append(elementValues, elValue)
		}
		return functions.State{
			Value: &values.ArrValue{elementValues},
			Stack: inputState.Stack,
		}
	}
	return outputShape, action, nil
}

///////////////////////////////////////////////////////////////////////////////

type CompositionExpression struct {
	Pos   lexer.Position
	Left  Expression
	Right Expression
}

func (x *CompositionExpression) Typecheck(inputShape functions.Shape, params []*functions.Parameter) (functions.Shape, functions.Action, error) {
	if len(params) > 0 {
		return zeroShape, nil, errors.E(
			errors.Kind(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	middleShape, lAction, err := x.Left.Typecheck(inputShape, nil)
	if err != nil {
		return zeroShape, nil, err
	}
	outputShape, rAction, err := x.Right.Typecheck(middleShape, nil)
	if err != nil {
		return zeroShape, nil, err
	}
	action := func(inputState functions.State, args []functions.Action) functions.State {
		middleState := lAction(inputState, nil)
		outputState := rAction(middleState, nil)
		return outputState
	}
	return outputShape, action, nil
}

///////////////////////////////////////////////////////////////////////////////

type CallExpression struct {
	Pos  lexer.Position
	Name string
	Args []Expression
}

func (x *CallExpression) Typecheck(inputShape functions.Shape, params []*functions.Parameter) (functions.Shape, functions.Action, error) {
	// Go down the function stack and find the function invoked by this
	// call
	stack := inputShape.FuncerStack
	for {
		// Reached bottom of stack without finding a matching function
		if stack == nil {
			return zeroShape, nil, errors.E(
				errors.Kind(errors.NoSuchFunction),
				errors.Pos(x.Pos),
				errors.InputType(inputShape.Type),
				errors.Name(x.Name),
				errors.NumParams(len(x.Args)+len(params)),
			)
		}
		// Try the funcer on top of the stack
		funcer := stack.Head
		funParams, funOutputType, funAction, ok := funcer(inputShape.Type, x.Name, len(x.Args)+len(params))
		if !ok {
			stack = stack.Tail
			continue
		}
		// Prepare action:
		action := funAction
		// Check function params filled by this call
		for i := 0; i < len(x.Args); i++ {
			param := funParams[i]
			argInputShape := functions.Shape{param.InputType, inputShape.FuncerStack}
			argOutputShape, argAction, err := x.Args[i].Typecheck(argInputShape, param.Params)
			if err != nil {
				return zeroShape, nil, err
			}
			if !param.OutputType.Subsumes(argOutputShape.Type) {
				return zeroShape, nil, errors.E(
					errors.Kind(errors.ArgHasWrongOutputType),
					errors.Pos(x.Pos),
					errors.ArgNum(i),
					errors.WantType(param.OutputType),
					errors.GotType(argOutputShape.Type),
				)
			}
			action = action.SetArg(argAction)
		}
		// Check function params *not* filled by this call (thus left
		// to function to call)
		for i := 0; i < len(params); i++ {
			if !params[i].Subsumes(funParams[len(x.Args)+i]) {
				return zeroShape, nil, errors.E(
					errors.Kind(errors.ParamDoesNotMatch),
					errors.Pos(x.Pos),
					errors.ParamNum(i),
					errors.WantParam(params[i]),
					errors.GotParam(funParams[len(x.Args)+i]),
				)
			}
		}
		// Return result
		return functions.Shape{funOutputType, inputShape.FuncerStack}, action, nil
	}
}

///////////////////////////////////////////////////////////////////////////////

type AssignmentExpression struct {
	Pos  lexer.Position
	Name string
}

func (x *AssignmentExpression) Typecheck(inputShape functions.Shape, params []*functions.Parameter) (functions.Shape, functions.Action, error) {
	if len(params) > 0 {
		return zeroShape, nil, errors.E(
			errors.Kind(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	var id interface{} = x
	varFuncer := func(gotInputType types.Type, gotName string, gotNumArgs int) ([]*functions.Parameter, types.Type, functions.Action, bool) {
		if gotName != x.Name {
			return nil, nil, nil, false
		}
		if gotNumArgs != 0 {
			return nil, nil, nil, false
		}
		varAction := func(inputState functions.State, args []functions.Action) functions.State {
			stack := inputState.Stack
			for stack != nil {
				if stack.Head.ID == id {
					return functions.State{
						Value: stack.Head.Action(functions.InitialState, nil).Value,
						Stack: inputState.Stack,
					}
				}
				stack = stack.Tail
			}
			panic(fmt.Sprintf("variable %s not found", x.Name))
		}
		return nil, inputShape.Type, varAction, true
	}
	outputShape := functions.Shape{inputShape.Type, inputShape.FuncerStack.Push(varFuncer)}
	action := func(inputState functions.State, args []functions.Action) functions.State {
		return functions.State{
			Value: inputState.Value,
			Stack: inputState.Stack.Push(functions.Variable{
				ID: id,
				Action: func(i functions.State, as []functions.Action) functions.State {
					return functions.State{
						Value: inputState.Value,
						Stack: i.Stack,
					}
				},
			}),
		}
	}
	return outputShape, action, nil
}

type valueStack struct {
	Head values.Value
	Tail *valueStack
}

func (s *valueStack) Push(element values.Value) *valueStack {
	return &valueStack{element, s}
}

///////////////////////////////////////////////////////////////////////////////

type DefinitionExpression struct {
	Pos        lexer.Position
	InputType  types.Type
	Name       string
	Params     []*functions.Parameter
	OutputType types.Type
	Body       Expression
}

func (x *DefinitionExpression) Typecheck(inputShape functions.Shape, params []*functions.Parameter) (functions.Shape, functions.Action, error) {
	// make sure we got no parameters
	if len(params) > 0 {
		return zeroShape, nil, errors.E(
			errors.Kind(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	// variables for body input stack, action (will be set at runtime)
	var bodyInputStackStub *functions.VariableStack
	var bodyAction functions.Action
	// make a funcer for the defined function, add it to the function stack
	funFuncer := func(gotInputType types.Type, gotName string, gotNumArgs int) ([]*functions.Parameter, types.Type, functions.Action, bool) {
		if !x.InputType.Subsumes(gotInputType) {
			return nil, nil, nil, false
		}
		if gotName != x.Name {
			return nil, nil, nil, false
		}
		if gotNumArgs != len(x.Params) {
			return nil, nil, nil, false
		}
		funAction := func(inputState functions.State, args []functions.Action) functions.State {
			// Bind parameters to arguments by adding corresponding
			// Variable objects to the body input state.
			bodyInputStack := bodyInputStackStub
			for i, param := range x.Params {
				var id interface{} = param
				bodyInputStack = bodyInputStack.Push(functions.Variable{
					ID:     id,
					Action: args[i],
				})
			}
			bodyInputState := functions.State{
				Value: inputState.Value,
				Stack: bodyInputStack,
			}
			bodyOutputState := bodyAction(bodyInputState, nil)
			return functions.State{
				Value: bodyOutputState.Value,
				Stack: inputState.Stack,
			}
		}
		return x.Params, x.OutputType, funAction, true
	}
	functionStack := inputShape.FuncerStack.Push(funFuncer)
	// add parameter funcers for use in the body
	bodyFuncerStack := functionStack
	for _, param := range x.Params {
		var id interface{} = param
		paramFuncer := func(gotInputType types.Type, gotName string, gotNumArgs int) ([]*functions.Parameter, types.Type, functions.Action, bool) {
			if !param.InputType.Subsumes(gotInputType) {
				return nil, nil, nil, false
			}
			if gotName != param.Name {
				return nil, nil, nil, false
			}
			if gotNumArgs != len(param.Params) {
				return nil, nil, nil, false
			}
			paramAction := func(inputState functions.State, args []functions.Action) functions.State {
				stack := inputState.Stack
				for stack != nil {
					if stack.Head.ID == id {
						return stack.Head.Action(inputState, args)
					}
					stack = stack.Tail
				}
				panic("action not found")
			}
			return param.Params, param.OutputType, paramAction, true
		}

		bodyFuncerStack = bodyFuncerStack.Push(paramFuncer)
	}
	// define body input context
	bodyInputShape := functions.Shape{
		Type:        x.InputType,
		FuncerStack: bodyFuncerStack,
	}
	// typecheck body (crucially, setting body action)
	bodyOutputShape, bodyAction, err := x.Body.Typecheck(bodyInputShape, nil)
	if err != nil {
		return zeroShape, nil, err
	}
	// check body output type
	if !x.OutputType.Subsumes(bodyOutputShape.Type) {
		return zeroShape, nil, errors.E(
			errors.Kind(errors.FunctionBodyHasWrongOutputType),
			errors.Pos(x.Pos),
			errors.WantType(x.OutputType),
			errors.GotType(bodyOutputShape.Type),
		)
	}
	// define output context
	outputShape := functions.Shape{
		Type:        inputShape.Type,
		FuncerStack: functionStack,
	}
	// define action (crucially, setting body input stack)
	action := func(inputState functions.State, args []functions.Action) functions.State {
		bodyInputStackStub = inputState.Stack
		return inputState
	}
	// return
	return outputShape, action, nil
}

///////////////////////////////////////////////////////////////////////////////

type ConditionalExpression struct {
	Pos             lexer.Position
	Condition       Expression
	Consequent      Expression
	ElifConditions  []Expression
	ElifConsequents []Expression
	Alternative     Expression
}

func (x *ConditionalExpression) Typecheck(inputShape functions.Shape, params []*functions.Parameter) (functions.Shape, functions.Action, error) {
	// make sure we got no parameters
	if len(params) > 0 {
		return zeroShape, nil, errors.E(
			errors.Kind(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	// typecheck condition
	conditionOutputShape, conditionAction, err := x.Condition.Typecheck(inputShape, nil)
	if err != nil {
		return zeroShape, nil, err
	}
	if !boolType.Subsumes(conditionOutputShape.Type) {
		return zeroShape, nil, errors.E(
			errors.Kind(errors.ConditionMustBeBool),
			errors.Pos(x.Pos),
			errors.WantType(boolType),
			errors.GotType(conditionOutputShape.Type),
		)
	}
	// context is the shared input context for all conditions and consequents.
	// Each condition may add to the FuncerStack. Type always stays the same.
	shape := functions.Shape{
		Type:        inputShape.Type,
		FuncerStack: conditionOutputShape.FuncerStack,
	}
	consequentOutputShape, consequentAction, err := x.Consequent.Typecheck(shape, nil)
	if err != nil {
		return zeroShape, nil, err
	}
	outputType := consequentOutputShape.Type
	elifConditionActions := make([]functions.Action, 0, len(x.ElifConditions))
	elifConsequentActions := make([]functions.Action, 0, len(x.ElifConsequents))
	for i := range x.ElifConditions {
		conditionOutputShape, elifConditionAction, err := x.ElifConditions[i].Typecheck(shape, nil)
		if err != nil {
			return zeroShape, nil, err
		}
		if !boolType.Subsumes(conditionOutputShape.Type) {
			return zeroShape, nil, errors.E(
				errors.Kind(errors.ConditionMustBeBool),
				errors.Pos(x.Pos),
				errors.WantType(boolType),
				errors.GotType(conditionOutputShape.Type),
			)
		}
		shape.FuncerStack = conditionOutputShape.FuncerStack
		elifConditionActions = append(elifConditionActions, elifConditionAction)
		consequentOutputShape, elifConsequentAction, err := x.ElifConsequents[i].Typecheck(shape, nil)
		if err != nil {
			return zeroShape, nil, err
		}
		elifConsequentActions = append(elifConsequentActions, elifConsequentAction)
		outputType = types.Disjoin(outputType, consequentOutputShape.Type)
	}
	alternativeOutputShape, alternativeAction, err := x.Alternative.Typecheck(shape, nil)
	if err != nil {
		return zeroShape, nil, err
	}
	outputType = types.Disjoin(outputType, alternativeOutputShape.Type)
	action := func(inputState functions.State, args []functions.Action) functions.State {
		conditionState := conditionAction(inputState, nil)
		boolConditionValue, _ := conditionState.Value.(*values.BoolValue)
		if boolConditionValue.Value {
			consequentInputState := functions.State{
				Value: inputState.Value,
				Stack: conditionState.Stack,
			}
			consequentOutputState := consequentAction(consequentInputState, nil)
			return functions.State{
				Value: consequentOutputState.Value,
				Stack: inputState.Stack,
			}
		}
		for i := range elifConditionActions {
			conditionState = elifConditionActions[i](inputState, nil)
			boolConditionValue, _ = conditionState.Value.(*values.BoolValue)
			if boolConditionValue.Value {
				consequentInputState := functions.State{
					Value: inputState.Value,
					Stack: conditionState.Stack,
				}
				consequentOutputState := elifConsequentActions[i](consequentInputState, nil)
				return functions.State{
					Value: consequentOutputState.Value,
					Stack: inputState.Stack,
				}
			}
		}
		alternativeOutputState := alternativeAction(inputState, nil)
		return functions.State{
			Value: alternativeOutputState.Value,
			Stack: inputState.Stack,
		}
	}
	outputShape := functions.Shape{
		Type:        outputType,
		FuncerStack: inputShape.FuncerStack,
	}
	return outputShape, action, nil
}

///////////////////////////////////////////////////////////////////////////////

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
	if !seqType.Subsumes(inputShape.Type) {
		return zeroShape, nil, errors.E(
			errors.Kind(errors.MappingRequiresSeqType),
			errors.Pos(x.Pos),
			errors.WantType(seqType),
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
		Type:        &types.SeqType{bodyOutputShape.Type},
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
			Value: &values.SeqValue{channel},
			Stack: inputState.Stack,
		}
	}
	return outputShape, action, nil
}

///////////////////////////////////////////////////////////////////////////////
