package errors

import (
	"fmt"
)

type ErrorCode int

const (
	Syntax ErrorCode = iota
	ParamsNotAllowed
	NoSuchFuncer
	ArgHasWrongOutputType
	ParamDoesNotMatch
	FunctionBodyHasWrongOutputType
	ConditionMustBeBool
	MappingRequiresArrType
	RestRequiresArrType
	ComposeWithVoid
	VoidProgram
	NonExhaustiveMatch
	ImpossibleMatch
	UnreachableElisClause
	UnreachableElseClause
	RegexpWantsString
	BadRegexp
	UnexpectedValue
	NoSuchProperty
	NoSuchIndex
	BadIndex
	NoGetterAllowed
)

func (code ErrorCode) String() string {
	switch code {
	case Syntax:
		return "Syntax"
	case ParamsNotAllowed:
		return "ParamsNotAllowed"
	case NoSuchFuncer:
		return "NoSuchFuncer"
	case ArgHasWrongOutputType:
		return "ArgHasWrongOutputType"
	case ParamDoesNotMatch:
		return "ParamDoesNotMatch"
	case FunctionBodyHasWrongOutputType:
		return "FunctionBodyHasWrongOutputType"
	case ConditionMustBeBool:
		return "ConditionMustBeBool"
	case MappingRequiresArrType:
		return "MappingRequiresArrType"
	case RestRequiresArrType:
		return "RestRequiresArrType"
	case ComposeWithVoid:
		return "ComposeWithVoid"
	case VoidProgram:
		return "VoidProgram"
	case NonExhaustiveMatch:
		return "NonExhaustiveMatch"
	case ImpossibleMatch:
		return "ImpossibleMatch"
	case UnreachableElisClause:
		return "UnreachableElisClause"
	case UnreachableElseClause:
		return "UnreachableElseClause"
	case RegexpWantsString:
		return "RegexpWantsString"
	case BadRegexp:
		return "BadRegexp"
	case UnexpectedValue:
		return "UnexpectedValue"
	case NoSuchProperty:
		return "NoSuchProperty"
	case NoSuchIndex:
		return "NoSuchIndex"
	case BadIndex:
		return "BadIndex"
	case NoGetterAllowed:
		return "NoGetterAllowed"
	default:
		return "Unknown"
	}
}

func (code ErrorCode) DefaultMessage() string {
	switch code {
	case Syntax:
		return "syntax error"
	case ParamsNotAllowed:
		return "This expression cannot be used as an argument here because it does not take parameters."
	case NoSuchFuncer:
		return "no such funcer"
	case ArgHasWrongOutputType:
		return "An argument has the wrong output type."
	case ParamDoesNotMatch:
		return "Cannot use this function here because one of its parameters does not match the expected parameter."
	case FunctionBodyHasWrongOutputType:
		return "The function body has the wrong output type."
	case ConditionMustBeBool:
		return "The condition must be boolean."
	case MappingRequiresArrType:
		return "The input to a mapping must be an array."
	case RestRequiresArrType:
		return "The rest of an array must itself be an array."
	case ComposeWithVoid:
		return "Cannot compose with a function whose return type is Void."
	case VoidProgram:
		return "Type of program cannot be Void."
	case NonExhaustiveMatch:
		return "Match is not exhaustive. Consider adding `elis` clauses and/or an `else` clause."
	case ImpossibleMatch:
		return "Impossible match. The pattern will never match the input type."
	case UnreachableElisClause:
		return "The `elis` clause is unreachable because the match is already exhaustive."
	case UnreachableElseClause:
		return "The `else` clause is unreachable because the match is already exhaustive."
	case RegexpWantsString:
		return "Regular expressions require Str as input type."
	case BadRegexp:
		return "error parsing regexp"
	case UnexpectedValue:
		return "Component got an unexpected input value."
	case NoSuchProperty:
		return "The object does not have this property."
	case NoSuchIndex:
		return "Array is not long enough."
	case BadIndex:
		return "Index must be a nonnegative integer."
	case NoGetterAllowed:
		return "A getter expression cannot be applied to this type."
	default:
		return "unknown error"
	}
}

func ParseCode(s string) (ErrorCode, error) {
	switch s {
	case "Syntax":
		return Syntax, nil
	case "ParamsNotAllowed":
		return ParamsNotAllowed, nil
	case "NoSuchFuncer":
		return NoSuchFuncer, nil
	case "ArgHasWrongOutputType":
		return ArgHasWrongOutputType, nil
	case "ParamDoesNotMatch":
		return ParamDoesNotMatch, nil
	case "FunctionBodyHasWrongOutputType":
		return FunctionBodyHasWrongOutputType, nil
	case "ConditionMustBeBool":
		return ConditionMustBeBool, nil
	case "MappingRequiresArrType":
		return MappingRequiresArrType, nil
	case "RestRequiresArrType":
		return RestRequiresArrType, nil
	case "ComposeWithVoid":
		return ComposeWithVoid, nil
	case "VoidProgram":
		return VoidProgram, nil
	case "NonExhaustiveMatch":
		return NonExhaustiveMatch, nil
	case "ImpossibleMatch":
		return ImpossibleMatch, nil
	case "UnreachableElisClause":
		return UnreachableElisClause, nil
	case "UnreachableElseClause":
		return UnreachableElseClause, nil
	case "RegexpWantsString":
		return RegexpWantsString, nil
	case "BadRegexp":
		return BadRegexp, nil
	case "UnexpectedValue":
		return UnexpectedValue, nil
	case "NoSuchProperty":
		return NoSuchProperty, nil
	case "NoSuchIndex":
		return NoSuchIndex, nil
	case "BadIndex":
		return BadIndex, nil
	case "NoGetterAllowed":
		return NoGetterAllowed, nil
	default:
		return 0, fmt.Errorf("invalid error code")
	}
}
