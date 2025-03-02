package docutil

import (
	"encoding/json"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/grammar"
	"github.com/texttheater/bach/interpreter"
)

func ParseError(input string) (error, error) {
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
		e.Kind = kind.(*errors.ErrorKind)
	}
	if code, ok := v["Code"]; ok {
		e.Code = code.(*errors.ErrorCode)
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
		*e.ArgNum = argNum.(int)
	}
	if numParams, ok := v["NumParams"]; ok {
		*e.NumParams = numParams.(int)
	}
	if paramNum, ok := v["ParamNum"]; ok {
		*e.ParamNum = paramNum.(int)
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
