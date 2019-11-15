package errors

type ErrorCode int

const (
	Syntax ErrorCode = iota
	ParamsNotAllowed
	NoSuchFunction
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
)

func (code ErrorCode) String() string {
	switch code {
	case Syntax:
		return "Syntax"
	case ParamsNotAllowed:
		return "ParamsNotAllowed"
	case NoSuchFunction:
		return "NoSuchFunction"
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
	case NoSuchFunction:
		return "no such function"
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
		return "Conditional got an unexpected input value."
	default:
		return "unknown error"
	}
}

func (code ErrorCode) Kind() string {
	switch code {
	case Syntax:
		return "Syntax error"
	case ParamsNotAllowed:
		return "Type error"
	case NoSuchFunction:
		return "Type error"
	case ArgHasWrongOutputType:
		return "Type error"
	case ParamDoesNotMatch:
		return "Type error"
	case FunctionBodyHasWrongOutputType:
		return "Type error"
	case ConditionMustBeBool:
		return "Type error"
	case MappingRequiresArrType:
		return "Type error"
	case RestRequiresArrType:
		return "Type error"
	case ComposeWithVoid:
		return "Type error"
	case VoidProgram:
		return "Type error"
	case NonExhaustiveMatch:
		return "Type error"
	case ImpossibleMatch:
		return "Type error"
	case UnreachableElisClause:
		return "Type error"
	case UnreachableElseClause:
		return "Type error"
	case RegexpWantsString:
		return "Type error"
	case BadRegexp:
		return "Syntax error"
	case UnexpectedValue:
		return "Value error"
	default:
		return "Unknown error"
	}
}
