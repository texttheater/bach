package tests

import (
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func SequenceTestCases() []TestCase {
	return []TestCase{
		{`for Seq<Num> def f Seq<Num> as =x x ok [1, 2, 3] f`, &types.SeqType{types.NumType{}}, values.ArrValue([]values.Value{values.NumValue(1), values.NumValue(2), values.NumValue(3)}), nil},
		{`[1, 2, 3] each *2 all arr`, &types.ArrType{types.NumType{}}, values.ArrValue([]values.Value{values.NumValue(2), values.NumValue(4), values.NumValue(6)}), nil},
		{`1 each *2 all`, nil, nil, errors.E(errors.Code(errors.MappingRequiresSeqType))},
	}
}
