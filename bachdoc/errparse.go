package main

import (
	"encoding/json"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/grammar"
	"github.com/texttheater/bach/interpreter"
)

func parseError(input string) (error, error) {
	var v map[string]any
	err := json.Unmarshal([]byte(input), &v)
	if err != nil {
		return nil, err
	}
	if len(v) == 0 {
		return nil, nil
	}
	var e errors.E
	if kind, ok := v["Kind"]; ok {
		kind, err := errors.ParseKind(kind.(string))
		if err != nil {
			return nil, err
		}
		e.Kind = &kind
	}
	if code, ok := v["Code"]; ok {
		code, err := errors.ParseCode(code.(string))
		if err != nil {
			return nil, err
		}
		e.Code = &code
	}
	if pos, ok := v["Pos"]; ok {
		pos := pos.(map[string]any)
		*e.Pos = lexer.Position{}
		if filename, ok := pos["Filename"]; ok {
			e.Pos.Filename = filename.(string)
		}
		if offset, ok := pos["Offset"]; ok {
			e.Pos.Offset = offset.(int)
		}
		if line, ok := pos["Line"]; ok {
			e.Pos.Line = line.(int)
		}
		if column, ok := pos["Column"]; ok {
			e.Pos.Column = column.(int)
		}
	}
	if message, ok := v["Message"]; ok {
		*e.Message = message.(string)
	}
	if wantType, ok := v["WantType"]; ok {
		wantType, err := grammar.ParseType(wantType.(string))
		if err != nil {
			return nil, err
		}
		e.WantType = wantType
	}
	if gotType, ok := v["GotType"]; ok {
		gotType, err := grammar.ParseType(gotType.(string))
		if err != nil {
			return nil, err
		}
		e.GotType = gotType
	}
	if gotValue, ok := v["GotValue"]; ok {
		_, v, err := interpreter.InterpretString(gotValue.(string))
		if err != nil {
			return nil, err
		}
		e.GotValue = v
	}
	if inputType, ok := v["InputType"]; ok {
		inputType, err := grammar.ParseType(inputType.(string))
		if err != nil {
			return nil, err
		}
		e.GotType = inputType
	}
	if name, ok := v["Name"]; ok {
		*e.Name = name.(string)
	}
	if argNum, ok := v["ArgNum"]; ok {
		argNum := int(argNum.(float64))
		e.ArgNum = &argNum
	}
	if numParams, ok := v["NumParams"]; ok {
		numParams := int(numParams.(float64))
		e.NumParams = &numParams
	}
	if paramNum, ok := v["ParamNum"]; ok {
		paramNum := int(paramNum.(float64))
		e.ParamNum = &paramNum
	}
	if wantParam, ok := v["WantParam"]; ok {
		wantParam, err := grammar.ParseParam(wantParam.(string))
		if err != nil {
			return nil, err
		}
		e.WantParam = wantParam
	}
	if gotParam, ok := v["GotParam"]; ok {
		gotParam, err := grammar.ParseParam(gotParam.(string))
		if err != nil {
			return nil, err
		}
		e.GotParam = gotParam
	}
	return &e, nil
}
