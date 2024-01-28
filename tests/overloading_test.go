package tests_test

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/tests"
	"github.com/texttheater/bach/types"
)

func TestOverloading(t *testing.T) {
	tests.TestProgram(`for Bool def f Num as 1 ok for Num def f Num as 2 ok true f`, types.Num{}, states.NumValue(1), nil, t)
	tests.TestProgram(`for Bool def f Num as 1 ok for Num def f Num as 2 ok 1 f`, types.Num{}, states.NumValue(2), nil, t)
}
