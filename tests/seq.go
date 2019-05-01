package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestSequences(t *testing.T) {
	TestProgram(`for Seq<Num> def f Seq<Num> as =x x ok [1, 2, 3] f`, &types.SeqType{types.NumType{}}, values.ArrValue([]values.Value{values.NumValue(1), values.NumValue(2), values.NumValue(3)}), nil, t)
	TestProgram(`[1, 2, 3] each *2 all arr`, &types.ArrType{types.NumType{}}, values.ArrValue([]values.Value{values.NumValue(2), values.NumValue(4), values.NumValue(6)}), nil, t)
	TestProgram(`1 each *2 all`, nil, nil, errors.E(errors.Code(errors.MappingRequiresSeqType)), t)
}
