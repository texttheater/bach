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
	MappingRequiresSeqType
	TailRequiresArrType
	ComposeWithVoid
	NonExhaustiveMatch
	ImpossibleMatch
	UnreachableElseClause
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
	case FunctionBodyHasWrongOutputType:
		return "FunctionBodyHasWrongOutputType"
	case ConditionMustBeBool:
		return "ConditionMustBeBool"
	case MappingRequiresSeqType:
		return "MappingRequiresSeqType"
	case TailRequiresArrType:
		return "TailRequiresArrType"
	case ComposeWithVoid:
		return "ComposeWithVoid"
	case NonExhaustiveMatch:
		return "NonExhaustiveMatch"
	case ImpossibleMatch:
		return "ImpossibleMatch"
	case UnreachableElseClause:
		return "UnreachableElseClause"
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
	case MappingRequiresSeqType:
		return "The input to a mapping must be a sequence."
	case TailRequiresArrType:
		return "The tail type of a Nearr type must be an array type."
	case ComposeWithVoid:
		return "Cannot compose with a function whose return type is Void."
	case NonExhaustiveMatch:
		return "Match is not exhaustive. Consider adding elis clauses and/or an else clause."
	case ImpossibleMatch:
		return "Impossible match. The pattern will never match the input type."
	case UnreachableElseClause:
		return "The `else` clause is unreachable because the match is already exhaustive."
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
	case MappingRequiresSeqType:
		return "Type error"
	case TailRequiresArrType:
		return "Type error"
	case ComposeWithVoid:
		return "Type error"
	case NonExhaustiveMatch:
		return "TypeError"
	case ImpossibleMatch:
		return "TypeError"
	case UnreachableElseClause:
		return "UnreachableElseClause"
	default:
		return "unknown error"
	}
}
